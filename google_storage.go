package main

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"io"
	"sync"
	"time"
)

type googleStorage struct {
	Bucket string
	*sync.RWMutex
}

func (g *googleStorage) getBucket() string {
	g.RLock()
	defer g.RUnlock()

	return g.Bucket
}

func (g *googleStorage) newClient(ctx context.Context) (*storage.Client, error) {
	return storage.NewClient(ctx)
}

func (g *googleStorage) upload(ctx context.Context, client *storage.Client, object string, data []byte) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	wc := client.Bucket(g.getBucket()).Object(object).NewWriter(ctx)

	if _, err := io.Copy(wc, bytes.NewReader(data)); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}
