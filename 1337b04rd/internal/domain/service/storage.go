package service

import (
	"context"
	"io"
)

type StorageService interface {
	UploadFile(ctx context.Context, bucketName string, objectName string, file io.Reader, fileSize int64, contentType string) (string, error)
}
