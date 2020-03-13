package storage

import (
	"context"
	"fmt"
	"sync"
)

const bucket = "ls-lint"

type Storage interface {
	InitClient(ctx context.Context) error
	CloseClient() error
	GetBucket() string
	Upload(ctx context.Context, object string, data []byte) error
}

func GetStorage(storage string) (Storage, error) {
	switch storage {
	case "google":
		return &Google{
			Bucket:  bucket,
			RWMutex: new(sync.RWMutex),
		}, nil
	}

	return nil, fmt.Errorf("storage %s not found", storage)
}
