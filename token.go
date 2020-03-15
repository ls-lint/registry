package main

import (
	"github.com/jinzhu/gorm"
	"sync"
)

type Token struct {
	gorm.Model
	Token *string `gorm:"unique;not null" json:"token" binding:"required"`
	Read  *bool   `gorm:"not null" json:"read" binding:"required"`
	Write *bool   `gorm:"not null" json:"write" binding:"required"`

	UserId *uint `gorm:"not null" json:"userId"`
	User   *User `json:"user"`

	*sync.RWMutex `json:"-"`
}

func (t *Token) init() {
	t.RWMutex = new(sync.RWMutex)
}

func (t *Token) getId() uint {
	t.RLock()
	defer t.RUnlock()

	return t.ID
}

func (t *Token) getToken() *string {
	t.RLock()
	defer t.RUnlock()

	return t.Token
}

func (t *Token) getUserId() *uint {
	t.RLock()
	defer t.RUnlock()

	return t.UserId
}

func (t *Token) setUserId(userId *uint) {
	t.Lock()
	defer t.Unlock()

	t.UserId = userId
}

func (t *Token) canWrite() bool {
	t.RLock()
	defer t.RUnlock()

	return *t.Write
}
