package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type api struct {
	database      *database
	googleStorage *googleStorage
	port          string
	mode          string
	*sync.RWMutex
}

func (api *api) getMode() string {
	api.RLock()
	defer api.RUnlock()

	return api.mode
}

func (api *api) getPort() string {
	api.RLock()
	defer api.RUnlock()

	return api.port
}

func (api *api) response(error, data interface{}) gin.H {
	return gin.H{
		"error": error,
		"data":  data,
	}
}

func (api *api) cors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")
	return cors.New(config)
}

func (api *api) startServer() error {
	gin.SetMode(api.getMode())
	r := gin.Default()
	r.Use(api.cors())

	apiGroup := r.Group("/api")
	{
		// health
		apiGroup.GET("/ok", api.ok)

		// publish
		apiGroup.POST("/publish", api.publish)
	}

	if os.Getenv("ENV") != "prod" {
		return r.Run(":" + api.getPort())
	}

	return r.RunTLS(":"+api.getPort(), os.Getenv("SSL_CERT"), os.Getenv("SSL_KEY"))
}

func (api *api) ok(c *gin.Context) {
	c.JSON(http.StatusOK, api.response(nil, "ok"))
}

func (api *api) getObject(user *User, _package *Package, release *Release) string {
	if release.Version != nil {
		return fmt.Sprintf("%s/%s/%s/%s", *user.Username, *_package.Name, *release.Version, ".ls-lint.yml")
	}

	return fmt.Sprintf("%s/%s/%s", *user.Username, *_package.Name, ".ls-lint.yml")
}

func (api *api) publish(c *gin.Context) {
	var (
		err     error
		publish *Publish
		client  *storage.Client
	)

	// context
	ctx := context.Background()

	// binding
	if err := c.Bind(&publish); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.response(err.Error(), nil))
		return
	}

	// login
	user, err := api.database.tokenLogin(publish.Token)

	if err != nil && err != gorm.ErrRecordNotFound {
		c.AbortWithStatusJSON(http.StatusUnauthorized, api.response(err.Error(), nil))
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.response(fmt.Errorf("authorization failed").Error(), nil))
		return
	}

	if !*user.Tokens[0].Write {
		c.AbortWithStatusJSON(http.StatusUnauthorized, api.response(fmt.Errorf("token has no write permissions").Error(), nil))
		return
	}

	// mimetype
	if !strings.HasPrefix(http.DetectContentType(publish.Data), "text/plain") {
		c.AbortWithStatusJSON(http.StatusForbidden, api.response(fmt.Errorf("invalid mimetype").Error(), nil))
		return
	}

	// publish package
	_package, err := api.database.publishPackage(publish, user)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.response(fmt.Errorf("publish package failed").Error(), nil))
		return
	}

	// publish release
	release, err := api.database.publishRelease(publish, _package)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.response(fmt.Errorf("publish release failed").Error(), nil))
		return
	}

	// google storage client
	if client, err = api.googleStorage.newClient(ctx); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.response(err.Error(), nil))
		return
	}

	// object
	object := api.getObject(user, _package, release)

	// upload
	if err := api.googleStorage.upload(ctx, client, object, publish.Data); err != nil {
		log.Fatal(err)
	}

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, api.response(nil, "ok"))
}
