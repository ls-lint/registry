package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func (api *api) _package(c *gin.Context) {
	var (
		err      error
		_package *Package
	)

	_package, err = api.database.publicPackage(c.Param("user"), c.Param("name"))

	if err != nil && err != gorm.ErrRecordNotFound {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.response(err.Error(), nil))
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, api.response(fmt.Errorf("package not found").Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.response(nil, _package))
}
