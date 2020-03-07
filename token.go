package main

import (
	"github.com/jinzhu/gorm"
	"sync"
)

type Token struct {
	gorm.Model
	Token *string `gorm:"not null" json:"token"`
	Read  *bool   `gorm:"not null" json:"read"`
	Write *bool   `gorm:"not null" json:"write"`

	UserId *uint `gorm:"not null" json:"userId"`
	User   *User `json:"user"`

	*sync.RWMutex `json:"-"`
}
