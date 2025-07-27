package config

import (
	"time"
)

type Config struct {
	Server   ServerConfig
	Storage  StorageConfig // minio
	Database DatabaseConfig
	Auth     AuthConfig
	App      AppConfig
}

type ServerConfig struct {
	Port         int           `env:"SERVER_PORT" default:"8080"`
	Host         string        `env:"SERVER_HOST" default:"0.0.0.0"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" default:"10s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" default:"10s"`
	BodyLimit    int           `env:"SERVER_BODY_LIMIT" default:"4194304"` // 4MB
}

type StorageConfig struct {
	Endpoint        string `env:"MINIO_ENDPOINT" default:"localhost:9000"`
	AccessKeyID     string `env:"MINIO_ACCESS_KEY" required:"true"`
	SecretAccessKey string `env:"MINIO_SECRET_KEY" required:"true"`
	BucketName      string `env:"MINIO_BUCKET" default:"quietstore-files"`
	UseSSL          bool   `env:"MINIO_USE_SSL" default:"false"`
	Region          string `env:"MINIO_REGION" default:"us-east-1"`
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" default:"localhost"`
	Port     int    `env:"DB_PORT" default:"5432"`
	User     string `env:"DB_USER" default:"postgres"`
	Password string `env:"DB_PASSWORD" required:"true"`
	Name     string `env:"DB_NAME" default:"quietstore"`
	SSLMode  string `env:"DB_SSL_MODE" default:"disable"`
}

type AuthConfig struct {
	JWTSecret     string        `env:"JWT_SECRET" required:"true"`
	JWTExpiration time.Duration `env:"JWT_EXPIRATION" default:"24h"`
	BCryptCost    int           `env:"BCRYPT_COST" default:"10"`
}

type AppConfig struct {
	Environment string `env:"APP_ENV" default:"development"`
	LogLevel    string `env:"LOG_LEVEL" default:"info"`
	MaxFileSize int64  `env:"MAX_FILE_SIZE" default:"104857600"` // 100MB
}
