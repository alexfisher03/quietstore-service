package config

import (
	"time"
)

type Config struct {
	Server  ServerConfig
	App     AppConfig
	Storage StorageConfig
}

type ServerConfig struct {
	Port         int           `env:"SERVER_PORT" default:"8080"`
	Host         string        `env:"SERVER_HOST" default:"0.0.0.0"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" default:"10000"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" default:"10000"`
	BodyLimit    int           `env:"SERVER_BODY_LIMIT" default:"41943040"`
}

type AppConfig struct {
	Environment string `env:"APP_ENVIRONMENT" default:"development"` // dev, test, prod
	LogLevel    string `env:"APP_LOG_LEVEL" default:"info"`
	MaxFileSize int64
}

type StorageConfig struct {
	BasePath string `env:"STORAGE_BASE_PATH" default:"./data"`
}
