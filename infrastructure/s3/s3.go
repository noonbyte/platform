package s3

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/noonbyte/platform/configs"
)

type S3Client interface {
	Client() *minio.Client
	Bucket() string
}

type s3Client struct {
	client *minio.Client
	bucket string
}

func New(cfg *configs.S3Configuration) (S3Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		),
		Region: cfg.Region,
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	if _, err := client.BucketExists(ctx, cfg.Bucket); err != nil {
		return nil, err
	}

	return &s3Client{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

func (c *s3Client) Client() *minio.Client {
	return c.client
}

func (c *s3Client) Bucket() string {
	return c.bucket
}
