package main

import (
	"github.com/ls-lint/registry/storage"
	"log"
	"os"
	"sync"
)

func main() {
	var err error

	// database
	database := &database{
		host:     os.Getenv("DB_HOST"),
		database: os.Getenv("DB_NAME"),
		port:     os.Getenv("DB_PORT"),
		username: os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASS"),
	}

	if database.connection, err = database.connect(); err != nil {
		log.Fatal(err)
	}

	database.migrate()

	// storage
	googleStorage, err := storage.GetStorage("google")

	if err != nil {
		log.Fatal(err)
	}

	// api
	api := &api{
		storage:         googleStorage,
		database:        database,
		port:            os.Getenv("API_PORT"),
		mode:            os.Getenv("API_MODE"),
		authIdentityKey: "id",
		RWMutex:         new(sync.RWMutex),
	}

	if err = api.startServer(); err != nil {
		log.Fatal(err)
	}
}
