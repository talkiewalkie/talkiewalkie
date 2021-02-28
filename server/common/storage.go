package common

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/satori/go.uuid"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
)

type StorageClient interface {
	Upload(ctx context.Context, blob io.Reader) (*uuid.UUID, error)
	Url(dest string) (string, error)
}

var _ StorageClient = GoogleStorage{}

type GoogleStorage struct {
	*storage.Client
	Cfg        *jwt.Config
	BucketName string
}

func initStorageClient(ctx context.Context) (StorageClient, error) {
	serviceAccountFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(serviceAccountFile))
	if err != nil {
		log.Fatalf("could not init storage client: %+v", err)
	}

	saKey, err := ioutil.ReadFile(serviceAccountFile)
	if err != nil {
		log.Fatalf("could not read service account file: %+v", err)
	}

	cfg, err := google.JWTConfigFromJSON(saKey)
	if err != nil {
		log.Fatalf("could not build jwt config from service account file: %+v", err)
	}

	g := GoogleStorage{Client: client, Cfg: cfg, BucketName: os.Getenv("BUCKET_NAME")}
	if g.BucketName == "" {
		return nil, fmt.Errorf("bad config: no bucket name")
	}
	return g, nil
}

func (g GoogleStorage) Upload(c context.Context, blob io.Reader) (*uuid.UUID, error) {
	uid := uuid.NewV4()
	remoteBlob := g.Bucket(g.BucketName).Object(uid.String())
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
	url, err := storage.SignedURL(g.BucketName, dest, &storage.SignedURLOptions{
		GoogleAccessID: g.Cfg.Email,
		PrivateKey:     g.Cfg.PrivateKey,
		Method:         http.MethodGet,
		Expires:        time.Now().Add(3 * time.Hour),
	})
	return url, err
}
