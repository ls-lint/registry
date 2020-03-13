package main

import (
	"github.com/jinzhu/gorm"
	"sync"
)

type User struct {
	gorm.Model
	Username *string `gorm:"not null" json:"username"`
	Password *string `gorm:"not null" json:"password"`
	Email    *string `gorm:"not null" json:"email"`

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

func (u *User) getTokens() []*Token {
	u.RLock()
	defer u.RUnlock()

	return u.Tokens
}
