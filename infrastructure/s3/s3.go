package s3

import (
	"context"
	"noonbyte/platform/configs"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Client interface {
	Client() *minio.Client
}

type s3Client struct {
	client *minio.Client
}

func New(cfg configs.S3Configuration) (S3Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	_, err = client.ListBuckets(context.Background())
	if err != nil {
		return nil, err
	}

	return &s3Client{
		client: client,
	}, nil
}

func (c *s3Client) Client() *minio.Client {
	return c.client
}
