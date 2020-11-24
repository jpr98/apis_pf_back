package datastore

import (
	"context"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
)

// StorageDatastore contains the information of Google Cloud Storage
type StorageDatastore struct {
	Client *storage.Client
	Bucket *storage.BucketHandle
	URL    string
}

// NewStorageDatastore creates a new StorageDatastore
func NewStorageDatastore(bucket string, log echo.Logger) (*StorageDatastore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	bucketHandle := client.Bucket(bucket)

	return &StorageDatastore{
		Client: client,
		Bucket: bucketHandle,
		URL:    "https://storage.googleapis.com/" + bucket,
	}, nil
}

// Upload uploads a file with a given name to GCP Storage
func (sd *StorageDatastore) Upload(name string, file io.Reader) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	wc := sd.Bucket.Object(name).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}
