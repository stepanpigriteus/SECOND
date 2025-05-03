package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Мокируем переменные окружения
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSLMODE", "disable")

	os.Setenv("S3_ENDPOINT", "s3.amazonaws.com")
	os.Setenv("S3_ACCESS_KEY", "access_key")
	os.Setenv("S3_SECRET_KEY", "secret_key")
	os.Setenv("S3_BUCKET", "my-bucket")
	os.Setenv("S3_REGION", "us-west-2")

	// Загружаем конфигурацию
	cfg := LoadConfig()

	// Проверяем, что конфигурация соответствует ожидаемым значениям
	if cfg.DB.Host != "localhost" {
		t.Errorf("Expected DB.Host to be 'localhost', got %s", cfg.DB.Host)
	}
	if cfg.DB.Port != "5432" {
		t.Errorf("Expected DB.Port to be '5432', got %s", cfg.DB.Port)
	}
	if cfg.DB.User != "user" {
		t.Errorf("Expected DB.User to be 'user', got %s", cfg.DB.User)
	}
	if cfg.DB.Password != "password" {
		t.Errorf("Expected DB.Password to be 'password', got %s", cfg.DB.Password)
	}
	if cfg.DB.Name != "testdb" {
		t.Errorf("Expected DB.Name to be 'testdb', got %s", cfg.DB.Name)
	}
	if cfg.DB.SSLMode != "disable" {
		t.Errorf("Expected DB.SSLMode to be 'disable', got %s", cfg.DB.SSLMode)
	}

	if cfg.S3.Endpoint != "s3.amazonaws.com" {
		t.Errorf("Expected S3.Endpoint to be 's3.amazonaws.com', got %s", cfg.S3.Endpoint)
	}
	if cfg.S3.AccessKey != "access_key" {
		t.Errorf("Expected S3.AccessKey to be 'access_key', got %s", cfg.S3.AccessKey)
	}
	if cfg.S3.SecretKey != "secret_key" {
		t.Errorf("Expected S3.SecretKey to be 'secret_key', got %s", cfg.S3.SecretKey)
	}
	if cfg.S3.Bucket != "my-bucket" {
		t.Errorf("Expected S3.Bucket to be 'my-bucket', got %s", cfg.S3.Bucket)
	}
	if cfg.S3.Region != "us-west-2" {
		t.Errorf("Expected S3.Region to be 'us-west-2', got %s", cfg.S3.Region)
	}
}

// Тестируем отсутствие переменной окружения
func TestLoadConfig_MissingEnv(t *testing.T) {
	// Мокируем переменные окружения
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSLMODE", "disable")

	// Убираем одну из переменных окружения, чтобы проверить, как это обработается
	os.Unsetenv("S3_BUCKET")

	// Загружаем конфигурацию
	cfg := LoadConfig()

	// Проверяем, что переменная окружения была пустой или нулевой
	if cfg.S3.Bucket != "" {
		t.Errorf("Expected S3.Bucket to be empty, got %s", cfg.S3.Bucket)
	}
}
