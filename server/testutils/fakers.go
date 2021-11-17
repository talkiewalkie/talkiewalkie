package testutils

import (
	"context"
	"io"

	uuid "github.com/satori/go.uuid"

	"github.com/talkiewalkie/talkiewalkie/clients"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type FakeStorageClient struct{}

func (f FakeStorageClient) AssetUrl(asset *models.Asset) (string, error) {
	return "https://some.fake.url/123", nil
}

func (f FakeStorageClient) DefaultBucket() string {
	return "test-bucket"
}

func (f FakeStorageClient) Download(blobName string, writer io.Writer) error {
	_, err := writer.Write([]byte("hello this is test content"))
	return err
}

func (f FakeStorageClient) Upload(ctx context.Context, blob io.ReadSeeker) (*uuid.UUID, error) {
	uid := uuid.NewV4()
	return &uid, nil
}

func (f FakeStorageClient) SignedUrl(bucket, blobName string) (string, error) {
	return "https://some.fake.url/123", nil
}

var _ clients.StorageClient = FakeStorageClient{}
