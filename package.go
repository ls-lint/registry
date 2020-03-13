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
	User   *User `json:"user"`

	Releases []*Release `json:"releases"`

	*sync.RWMutex `json:"-"`
}

func createPackage(publish *Publish, user *User) *Package {
	return &Package{
		Name:    &publish.Package,
		Public:  &publish.Public,
		User:    user,
		RWMutex: new(sync.RWMutex),
	}
}

func (p *Package) getId() *uint {
	p.RLock()
	defer p.RUnlock()

	return &p.ID
}

func (p *Package) getName() *string {
	p.RLock()
	defer p.RUnlock()

	return p.Name
}

func (p *Package) getUser() *User {
	p.Lock()
	defer p.Unlock()

	return p.User
}
