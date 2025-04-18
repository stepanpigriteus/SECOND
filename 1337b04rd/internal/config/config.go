package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	S3 struct {
		Endpoint  string
		AccessKey string
		SecretKey string
		Bucket    string
		Region    string
	}
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{}
	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port = os.Getenv("DB_PORT")
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	cfg.DB.Name = os.Getenv("DB_NAME")
	cfg.DB.SSLMode = os.Getenv("DB_SSLMODE")

	cfg.S3.Endpoint = os.Getenv("S3_ENDPOINT")
	cfg.S3.AccessKey = os.Getenv("S3_ACCESS_KEY")
	cfg.S3.SecretKey = os.Getenv("S3_SECRET_KEY")
	cfg.S3.Bucket = os.Getenv("S3_BUCKET")
	cfg.S3.Region = os.Getenv("S3_REGION")
	
	return cfg
}
