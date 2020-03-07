package main

import (
	"github.com/jinzhu/gorm"
	"sync"
)

type Package struct {
	gorm.Model
	Name   *string `gorm:"not null" json:"name"`
	Public *bool   `gorm:"not null" json:"public"`

	UserId *uint `gorm:"not null" json:"userId"`
	User *User `json:"user"`

	Releases []*Release `json:"releases"`

	*sync.RWMutex `json:"-"`
}
