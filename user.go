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
