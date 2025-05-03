package api

import (
	"bytes"
	"context"
	"testing"
)

func TestNewMinioStorage(t *testing.T) {
	// Тестируем создание нового MinioStorage
	storage, err := NewMinioStorage("localhost:9000", "minioadmin_new", "minioadmin_new_password", false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if storage == nil {
		t.Fatalf("Expected a MinioStorage instance, got nil")
	}
}

func TestUploadFile(t *testing.T) {
	// Создаем фейковый файл для загрузки
	content := []byte("this is a test file")
	file := bytes.NewReader(content)
	fileSize := int64(len(content))
	contentType := "text/plain"

	// Мокируем MinioStorage
	storage, err := NewMinioStorage("localhost:9000", "minioadmin_new", "minioadmin_new_password", false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Создаем фейковый bucket
	bucketName := "posts"
	objectName := "testfile.txt"

	// Выполняем загрузку
	ctx := context.Background()
	url, err := storage.UploadFile(ctx, bucketName, objectName, file, fileSize, contentType)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Проверяем, что URL выглядит как ожидаемый
	expectedURL := "https://localhost:9000/posts/testfile.txt"
	if url != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, url)
	}
}

func TestUploadFile_BucketAlreadyExists(t *testing.T) {
	// Создаем фейковый файл для загрузки
	content := []byte("this is a test file")
	file := bytes.NewReader(content)
	fileSize := int64(len(content))
	contentType := "text/plain"

	// Мокируем MinioStorage
	storage, err := NewMinioStorage("localhost:9000", "minioadmin_new", "minioadmin_new_password", false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Пытаемся загрузить файл в уже существующий бакет
	bucketName := "posts" // Бакет уже существует
	objectName := "existingfile.txt"
	ctx := context.Background()

	// Загружаем файл
	_, err = storage.UploadFile(ctx, bucketName, objectName, file, fileSize, contentType)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestUploadFile_BucketNotExists(t *testing.T) {
	// Создаем фейковый файл для загрузки
	content := []byte("this is a test file")
	file := bytes.NewReader(content)
	fileSize := int64(len(content))
	contentType := "text/plain"

	// Мокируем MinioStorage с несуществующим bucket
	storage, err := NewMinioStorage("localhost:9000", "minioadmin_new", "minioadmin_new_password", false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Попытка загрузить в несуществующий бакет
	bucketName := "nonexistentbucket"
	objectName := "testfile.txt"
	ctx := context.Background()

	_, err = storage.UploadFile(ctx, bucketName, objectName, file, fileSize, contentType)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestNewMinioStorage_BadCredentials(t *testing.T) {
	// Тестируем ошибку при передаче неверных данных
	storage, err := NewMinioStorage("localhost:9000", "wrongAccessKey", "wrongSecretKey", false)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	if storage != nil {
		t.Fatalf("Expected nil storage, got %v", storage)
	}
}

func TestNewMinioStorage_BadEndpoint(t *testing.T) {
	// Тестируем ошибку при передаче неправильного endpoint
	storage, err := NewMinioStorage("wrong-endpoint:9000", "minioadmin_new", "minioadmin_new_password", false)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	if storage != nil {
		t.Fatalf("Expected nil storage, got %v", storage)
	}
}

func TestUploadFile_WithDifferentContentTypes(t *testing.T) {
	// Тестируем загрузку с различными типами контента
	types := []string{"image/png", "application/pdf", "text/html"}
	for _, contentType := range types {
		t.Run(contentType, func(t *testing.T) {
			content := []byte("this is a test file for " + contentType)
			file := bytes.NewReader(content)
			fileSize := int64(len(content))

			// Мокируем MinioStorage
			storage, err := NewMinioStorage("localhost:9000", "minioadmin_new", "minioadmin_new_password", false)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			bucketName := "posts"
			objectName := "testfile.txt"
			ctx := context.Background()

			url, err := storage.UploadFile(ctx, bucketName, objectName, file, fileSize, contentType)
			if err != nil {
				t.Fatalf("Expected no error for content type %s, got %v", contentType, err)
			}

			// Проверяем URL
			expectedURL := "https://localhost:9000/posts/testfile.txt"
			if url != expectedURL {
				t.Errorf("Expected URL %s for content type %s, got %s", expectedURL, contentType, url)
			}
		})
	}
}
