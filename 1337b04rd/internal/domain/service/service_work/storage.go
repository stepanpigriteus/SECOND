package servicework

import (
	apiport "a1337b04rd/internal/domain/port/api"
	"context"
	"io"
)

type FileStorageService struct {
	storage apiport.FileStorage
}

func NewFileStorageService(storage apiport.FileStorage) *FileStorageService {
	return &FileStorageService{storage: storage}
}

func (s *FileStorageService) UploadFile(ctx context.Context, objectName string, file io.Reader, fileSize int64, contentType string) (string, error) {
	return s.storage.UploadFile(ctx, objectName, file, fileSize, contentType)
}
