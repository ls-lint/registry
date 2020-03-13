package storage

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"io"
	"sync"
	"time"
)

type Google struct {
	Client *storage.Client
	Bucket string
	*sync.RWMutex
}

func (g *Google) getClient() *storage.Client {
	g.RLock()
	defer g.RUnlock()

	return g.Client
}

func (g *Google) InitClient(ctx context.Context) (err error) {
	g.Lock()
	defer g.Unlock()

	g.Client, err = storage.NewClient(ctx)
	return err
}

func (g *Google) CloseClient() error {
	g.RLock()
	defer g.RUnlock()

	return g.Client.Close()
}

func (g *Google) GetBucket() string {
	g.RLock()
	defer g.RUnlock()

	return g.Bucket
}

func (g *Google) Upload(ctx context.Context, object string, data []byte) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	wc := g.getClient().Bucket(g.GetBucket()).Object(object).NewWriter(ctx)

	if _, err := io.Copy(wc, bytes.NewReader(data)); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}
