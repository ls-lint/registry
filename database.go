package main

import (
	"crypto/md5"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strings"
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

func (database *database) tokenLogin(token string) (*User, error) {
	var user User

	return &user, database.connection.
		Preload("Tokens", "token = ?", token).
		Joins("JOIN tokens ON users.id=tokens.id AND tokens.token = ?", token).
		First(&user).
		Error
}

func (database *database) publishPackage(publish *Publish, user *User) (*Package, error) {
	publish.Package = strings.ToLower(publish.Package)

	_package := &Package{
		Name:   &publish.Package,
		Public: &publish.Public,
		User:   user,
	}

	return _package, database.connection.
		FirstOrCreate(_package, &Package{
			Name:   _package.Name,
			UserId: &user.ID,
		}).
		Error
}

func (database *database) publishRelease(publish *Publish, _package *Package) (*Release, error) {
	hash := fmt.Sprintf("%x", md5.Sum(publish.Data))

	release := &Release{
		Version: &publish.Version,
		Hash:    &hash,
		Package: _package,
	}

	return release, database.connection.
		FirstOrCreate(release, &Release{
			Version:   &publish.Version,
			PackageId: &_package.ID,
		}).
		Error
}
