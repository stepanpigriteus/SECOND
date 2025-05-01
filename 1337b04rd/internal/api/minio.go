package api

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStorage(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*MinioStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// Ensure the bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &MinioStorage{client: client, bucketName: bucketName}, nil
}

func (s *MinioStorage) UploadFile(ctx context.Context, objectName string, file io.Reader, fileSize int64, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucketName, objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	// Return the URL (you may adjust this)
	url := fmt.Sprintf("https://%s/%s/%s", s.client.EndpointURL().Host, s.bucketName, objectName)
	return url, nil
}
