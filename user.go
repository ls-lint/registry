package main

import (
	"github.com/jinzhu/gorm"
	"sync"
)

type User struct {
	gorm.Model
	Username *string `gorm:"unique;not null" json:"username" binding:"required"`
	Password *string `gorm:"not null" json:"password" binding:"required"`
	Email    *string `gorm:"unique;not null" json:"email" binding:"required"`

	Packages []*Package `json:"packages"`
	Tokens   []*Token   `json:"tokens"`

	*sync.RWMutex `json:"-"`
}

func (u *User) init() {
	u.RWMutex = new(sync.RWMutex)

	for _, token := range u.Tokens {
		token.init()
	}
}

func (u *User) getId() *uint {
	u.RLock()
	defer u.RUnlock()

	return &u.ID
}

func (u *User) getUsername() *string {
	u.RLock()
	defer u.RUnlock()

	return u.Username
}

func (u *User) getPassword() *string {
	u.RLock()
	defer u.RUnlock()

	return u.Password
}

func (u *User) getTokens() []*Token {
	u.RLock()
	defer u.RUnlock()

	return u.Tokens
}
