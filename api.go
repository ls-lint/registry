package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ls-lint/registry/storage"
	"net/http"
	"os"
	"sync"
)

type api struct {
	database        *database
	storage         storage.Storage
	port            string
	mode            string
	authIdentityKey string
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

func (api *api) getAuthIdentityKey() string {
	api.RLock()
	defer api.RUnlock()

	return api.authIdentityKey
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
	r.Use(gin.Recovery())

	// auth
	authMiddleware, err := api.authMiddleware()

	if err != nil {
		return err
	}

	apiGroup := r.Group("/api")
	{
		// health
		apiGroup.GET("/ok", api.ok)

		// auth
		apiGroup.POST("/login", authMiddleware.LoginHandler)
		apiGroup.POST("/register", api.register)

		// publish
		apiGroup.POST("/publish", api.publish)

		// package
		apiGroup.GET("/package/:user/:name", api._package)

		// auth
		authGroup := apiGroup.Group("/auth")
		authGroup.Use(authMiddleware.MiddlewareFunc())
		{
			// auth
			authGroup.GET("/refresh_token", authMiddleware.RefreshHandler)
			authGroup.GET("/logout", authMiddleware.LogoutHandler)

			// token
			authGroup.GET("/token", api.tokens)
			authGroup.POST("/token", api.createToken)
			authGroup.DELETE("/token", api.deleteToken)

			// package
			authGroup.GET("/package/:user/:name", api._package)
		}
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
	if release.getTag() != nil {
		return fmt.Sprintf("%s/%s/%s/%s", *user.getUsername(), *_package.getName(), *release.getTag(), ".ls-lint.yml")
	}

	return fmt.Sprintf("%s/%s/%s", *user.getUsername(), *_package.getName(), ".ls-lint.yml")
}
