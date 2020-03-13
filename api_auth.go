package main

import (
	"github.com/appleboy/gin-jwt/v2"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func (api *api) authMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "auth",
		Key:         []byte(os.Getenv("JWT_KEY")),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: api.getAuthIdentityKey(),
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				v.init()

				return jwt.MapClaims{
					api.getAuthIdentityKey(): *v.getUsername(),
				}
			}

			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			spew.Dump(claims)
			username := claims[api.getAuthIdentityKey()].(string)

			return &User{
				Username: &username,
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			return api.database.login(c.PostForm("username"), c.PostForm("password"))
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, api.response(message, nil))
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}

func (api *api) register(c *gin.Context) {
	var user User

	// init
	user.init()

	// bind json
	if err := c.Bind(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.response(err.Error(), nil))
		return
	}

	// crypt password
	bcrypted, err := bcrypt.GenerateFromPassword([]byte(*user.Password), 14)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.response(err.Error(), nil))
		return
	}

	// update user password
	bcryptedString := string(bcrypted)
	user.Password = &bcryptedString

	// create user
	if err := api.database.connection.Create(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.response(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, api.response(nil, user))
}
