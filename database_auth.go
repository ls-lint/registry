package main

import (
	"github.com/appleboy/gin-jwt/v2"
	"golang.org/x/crypto/bcrypt"
)

func (database *database) tokenLogin(token string) (*User, error) {
	var user User

	return &user, database.connection.
		Preload("Tokens", "token = ?", token).
		Joins("JOIN tokens ON users.id=tokens.id AND tokens.token = ?", token).
		First(&user).
		Error
}

func (database *database) login(username string, password string) (*User, error) {
	var user User

	// init
	user.init()

	// query user
	if err := database.connection.Where("username = ?", username).Find(&user).Error; err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	// check password
	if bcrypt.CompareHashAndPassword([]byte(*user.getPassword()), []byte(password)) != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	return &user, nil
}
