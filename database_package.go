package main

func (database *database) publicPackage(user string, name string) (*Package, error) {
	var _package Package

	return &_package, database.connection.
		Preload("Releases").
		Where("users.username = ? AND packages.name = ? AND packages.public = ?", user, name, true).
		Joins("JOIN users ON users.id = packages.user_id").
		Find(&_package).
		Error
}
