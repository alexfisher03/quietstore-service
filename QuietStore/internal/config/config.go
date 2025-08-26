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
	Environment         string `env:"APP_ENVIRONMENT" default:"development"`
	LogLevel            string `env:"APP_LOG_LEVEL" default:"info"`
	MaxFileSize         int64
	RateLimitAuthMax    int           `env:"RATE_LIMIT_AUTH_MAX" default:"5"`
	RateLimitAuthExpire time.Duration `env:"RATE_LIMIT_AUTH_EXPIRATION" default:"60"`
	RateLimitUserMax    int           `env:"RATE_LIMIT_USER_MAX" default:"3"`
	RateLimitUserExpire time.Duration `env:"RATE_LIMIT_USER_EXPIRATION" default:"60"`
	RateLimitFileMax    int           `env:"RATE_LIMIT_FILE_MAX" default:"15"`
	RateLimitFileExpire time.Duration `env:"RATE_LIMIT_FILE_EXPIRATION" default:"60"`
}

type StorageConfig struct {
	BasePath string `env:"STORAGE_BASE_PATH" default:"./data"`
}
