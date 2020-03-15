package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (api *api) tokens(c *gin.Context) {
	var (
		err    error
		tokens []*Token
		userId = api.getUserId(c)
	)

	if tokens, err = api.database.tokens(userId); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.response(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.response(nil, tokens))
}

func (api *api) createToken(c *gin.Context) {
	var (
		err    error
		token  *Token
		userId = api.getUserId(c)
	)

	// bind
	if err := c.Bind(&token); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.response(err.Error(), nil))
		return
	}

	// init token
	token.init()
	token.setUserId(&userId)

	// create at database
	if token, err = api.database.createToken(token); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.response(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.response(nil, token))
}

func (api *api) deleteToken(c *gin.Context) {
	var (
		err    error
		token  *Token
		userId = api.getUserId(c)
	)

	// bind
	if err := c.Bind(&token); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.response(err.Error(), nil))
		return
	}

	// init token
	token.init()

	// verify token
	if *token.getUserId() != userId{
		c.AbortWithStatusJSON(http.StatusUnauthorized, api.response("authorization failed", nil))
		return
	}

	// delete at database
	if err = api.database.deleteToken(token); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.response(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.response(nil, token))
}
