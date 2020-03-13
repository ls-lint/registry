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

func (t *Token) init() {
	t.RWMutex = new(sync.RWMutex)
}

func (t *Token) getToken() *string {
	t.RLock()
	defer t.RUnlock()

	return t.Token
}

func (t *Token) canRead() bool {
	t.RLock()
	defer t.RUnlock()

	return *t.Read == true
}

func (t *Token) canWrite() bool {
	t.RLock()
	defer t.RUnlock()

	return *t.Write == true
}
