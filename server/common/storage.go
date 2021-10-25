package common

import (
	"context"
	"fmt"
	"github.com/talkiewalkie/talkiewalkie/models"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/satori/go.uuid"
)

type StorageClient interface {
	Upload(ctx context.Context, blob io.ReadSeeker) (*uuid.UUID, error)
	Download(blobName string, writer io.Writer) error
	SignedUrl(bucket, blobName string) (string, error)
	AssetUrl(asset *models.Asset) (string, error)
	DefaultBucket() string
}

var _ StorageClient = GoogleStorage{}

type GoogleStorage struct {
	*storage.Client
	BucketName string
}

func (g GoogleStorage) AssetUrl(asset *models.Asset) (string, error) {
	// TODO: coming back to this, it feels like superfluous complexity, should always precise the bucket name and not
	// 		 rely on inference when nil.
	if asset.Bucket.Valid {
		return g.SignedUrl(asset.Bucket.String, asset.BlobName.String)
	} else {
		return g.SignedUrl(g.DefaultBucket(), asset.UUID.String())
	}
}

func initStorageClient(ctx context.Context) (StorageClient, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("could not init storage client: %+v", err)
	}

	g := GoogleStorage{Client: client, BucketName: os.Getenv("BUCKET_NAME")}
	if g.BucketName == "" {
		return nil, fmt.Errorf("bad config: no bucket name")
	}
	return g, nil
}

func (g GoogleStorage) Upload(c context.Context, blob io.ReadSeeker) (*uuid.UUID, error) {
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

func (g GoogleStorage) SignedUrl(bucket, blobName string) (string, error) {
	url, err := storage.SignedURL(bucket, blobName, nil)
	return url, err
}

func (g GoogleStorage) Download(blobName string, writer io.Writer) error {
	blob := g.Bucket(g.BucketName).Object(blobName)

	reader, err := blob.NewReader(context.Background())
	if err != nil {
		return err
	}

	if _, err = io.Copy(writer, reader); err != nil {
		return err
	}

	return nil
}

func (g GoogleStorage) DefaultBucket() string {
	return g.BucketName
}

type S3Storage struct {
	*s3.S3
	bucketName string
	sess       *session.Session
}

func (s S3Storage) Upload(ctx context.Context, blob io.ReadSeeker) (*uuid.UUID, error) {
	uid := uuid.NewV4()
	_, err := s.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Body:   blob,
		Bucket: aws.String(s.DefaultBucket()),
		Key:    aws.String(uid.String()),
	})
	return &uid, err
}

func (s S3Storage) Download(blobName string, writer io.Writer) error {
	o, err := s.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.DefaultBucket()),
		Key:    aws.String(blobName),
	})
	if err != nil {
		return nil
	}

	if _, err = io.Copy(writer, o.Body); err != nil {
		return err
	}
	return nil
}

func (s S3Storage) SignedUrl(bucket, blobName string) (string, error) {
	req, _ := s.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(blobName),
	})

	urlStr, err := req.Presign(10 * time.Minute)
	if err != nil {
		return "", nil
	}

	return urlStr, nil
}

func (s S3Storage) AssetUrl(asset *models.Asset) (string, error) {
	if asset.Bucket.Valid {
		return s.SignedUrl(asset.Bucket.String, asset.BlobName.String)
	} else {
		return s.SignedUrl(s.DefaultBucket(), asset.UUID.String())
	}
}

func (s S3Storage) DefaultBucket() string {
	return s.bucketName
}

var _ StorageClient = S3Storage{}

func NewS3Storage() (*S3Storage, error) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	bn := os.Getenv("BUCKET_NAME")
	return &S3Storage{S3: svc, sess: sess, bucketName: bn}, nil
}
