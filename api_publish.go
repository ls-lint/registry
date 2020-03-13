package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"
)

func (api *api) publish(c *gin.Context) {
	var (
		err     error
		publish *Publish
	)

	// context
	ctx := context.Background()

	// binding
	if err := c.Bind(&publish); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.response(err.Error(), nil))
		return
	}

	// login
	user, err := api.database.getUserByToken(publish.Token)

	if err != nil && err != gorm.ErrRecordNotFound {
		c.AbortWithStatusJSON(http.StatusUnauthorized, api.response(err.Error(), nil))
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.response(fmt.Errorf("authorization failed").Error(), nil))
		return
	}

	// init user
	user.init()

	// validate token permissions
	if !user.getTokens()[0].canWrite() {
		c.AbortWithStatusJSON(http.StatusUnauthorized, api.response(fmt.Errorf("token has no write permissions").Error(), nil))
		return
	}

	// mimetype
	if !strings.HasPrefix(http.DetectContentType(publish.Data), "text/plain") {
		c.AbortWithStatusJSON(http.StatusForbidden, api.response(fmt.Errorf("invalid mimetype").Error(), nil))
		return
	}

	// create package
	_package := createPackage(publish, user)

	// create release
	release := createRelease(publish, _package)

	// client
	if err := api.storage.InitClient(ctx); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.response(err.Error(), nil))
		return
	}

	// upload
	if err := api.storage.Upload(ctx, api.getObject(user, _package, release), publish.Data); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.response(err.Error(), nil))
		return
	}

	// close client
	if err := api.storage.CloseClient(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.response(err.Error(), nil))
		return
	}

	// publish package
	_package, err = api.database.publishPackage(_package)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.response(fmt.Errorf("publish package failed").Error(), nil))
		return
	}

	// publish release
	release, err = api.database.publishRelease(release)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.response(fmt.Errorf("publish release failed").Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.response(nil, "ok"))
}
