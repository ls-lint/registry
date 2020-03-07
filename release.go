package main

import (
	"github.com/jinzhu/gorm"
	"sync"
)

type Release struct {
	gorm.Model
	Version *string `json:"version"`
	Hash    *string `json:"hash"`

	PackageId *uint    `gorm:"not null" json:"packageId"`
	Package   *Package `json:"packages"`

	*sync.RWMutex `json:"-"`
}
