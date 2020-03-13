package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type database struct {
	connection *gorm.DB
	username   string
	password   string
	database   string
	host       string
	port       string
}

func (database *database) getConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=UTC",
		database.username,
		database.password,
		database.host,
		database.port,
		database.database,
	)
}

func (database *database) connect() (*gorm.DB, error) {
	return gorm.Open("mysql", database.getConnectionString())
}

func (database *database) migrate() {
	database.connection.AutoMigrate(
		new(User),
		new(Package),
		new(Release),
		new(Token),
	)
}

func (database *database) publishPackage(_package *Package) error {
	return database.connection.
		FirstOrCreate(_package, &Package{
			Name:   _package.getName(),
			UserId: _package.getUser().getId(),
		}).
		Error
}

func (database *database) publishRelease(release *Release) error {
	return database.connection.
		FirstOrCreate(release, &Release{
			Tag:       release.getTag(),
			PackageId: release.getPackage().getId(),
		}).
		Error
}
