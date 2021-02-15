package common

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/docker/distribution/uuid"
	"google.golang.org/api/option"
)

type StorageClient interface {
	Upload(ctx context.Context, blob io.Reader) (*uuid.UUID, error)
	Url(dest string) (string, error)
}

var _ StorageClient = GoogleStorage{}

type GoogleStorage struct {
	*storage.Client
	bucket string
}

func initStorageClient(ctx context.Context) (StorageClient, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
	if err != nil {
		log.Fatal(err)
	}

	var g GoogleStorage
	g.Client = client
	g.bucket = os.Getenv("BUCKET_NAME")
	if g.bucket == "" {
		return nil, fmt.Errorf("bad config: no bucket name")
	}
	return g, nil
}

func (g GoogleStorage) Upload(c context.Context, blob io.Reader) (*uuid.UUID, error) {
	uid := uuid.Generate()
	remoteBlob := g.Bucket(g.bucket).Object(uid.String())
	wc := remoteBlob.NewWriter(c)
	if _, err := io.Copy(wc, blob); err != nil {
		return nil, err
	}
	if err := wc.Close(); err != nil {
		return nil, err
	}
	return &uid, nil
}

func (g GoogleStorage) Url(dest string) (string, error) {
	url, err := storage.SignedURL(g.bucket, dest, &storage.SignedURLOptions{})
	return url, err
}
