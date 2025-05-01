package apiport

import (
	"context"
	"io"
)

type FileStorage interface {
	UploadFile(ctx context.Context, objectName string, file io.Reader, fileSize int64, contentType string) (string, error)
}
