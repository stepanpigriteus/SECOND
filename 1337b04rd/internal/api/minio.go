package minioadapter

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	minio  *minio.Client
	bucket string
}

func NewClient(endpoint, accessKey, secretKey string, useSSL bool, bucket string) (*Client, error) {
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := cli.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := cli.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &Client{minio: cli, bucket: bucket}, nil
}

// Реализация интерфейса портов
func (c *Client) UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	objectName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
	uploadInfo, err := c.minio.PutObject(ctx, c.bucket, objectName, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", err
	}

	log.Println("Uploaded:", uploadInfo)
	return fmt.Sprintf("http://%s/%s/%s", c.minio.EndpointURL().Host, c.bucket, objectName), nil
}
