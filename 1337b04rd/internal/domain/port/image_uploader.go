package port

import (
	"context"
	"mime/multipart"
)

type ImageUploader interface {
	UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
}
