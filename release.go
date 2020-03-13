package main

import (
	"crypto/md5"
	"fmt"
	"github.com/jinzhu/gorm"
	"sync"
)

type Release struct {
	gorm.Model
	Tag  *string `json:"tag"`
	Hash *string `json:"hash"`

	PackageId *uint    `gorm:"not null" json:"packageId"`
	Package   *Package `json:"packages"`

	*sync.RWMutex `json:"-"`
}

func createRelease(publish *Publish, _package *Package) *Release {
	hash := fmt.Sprintf("%x", md5.Sum(publish.Data))

	return &Release{
		Tag:     &publish.Tag,
		Hash:    &hash,
		Package: _package,
		RWMutex: new(sync.RWMutex),
	}
}

func (r *Release) getPackage() *Package {
	r.RLock()
	defer r.RUnlock()

	return r.Package
}

func (r *Release) getTag() *string {
	r.RLock()
	defer r.RUnlock()

	return r.Tag
}

func (r *Release) getHash() *string {
	r.RLock()
	defer r.RUnlock()

	return r.Hash
}
