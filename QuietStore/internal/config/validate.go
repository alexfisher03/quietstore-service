package config

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

func (config *Config) Validate() error {
	var errs []string

	if config.Server.Port < 1 || config.Server.Port > 65535 {
		errs = append(errs, "server port must be between 1 and 65535")
	}

	if config.Storage.AccessKeyID == "" {
		errs = append(errs, "minio access key id must not be empty you fucker")
	}

	if config.Storage.SecretAccessKey == "" {
		errs = append(errs, "minio secret access key must not be empty you fucker")
	}

	if _, _, err := net.SplitHostPort(config.Storage.Endpoint); err != nil {
		errs = append(errs, "minio endpoint must be a valid host:port format")
	}

	if len(config.Auth.JWTSecret) < 32 {
		errs = append(errs, "JWT secret must be at least 32 characters")
	}

	if config.Auth.BCryptCost < 10 || config.Auth.BCryptCost > 31 {
		errs = append(errs, "BCrypt cost must be between 10 and 31")
	}

	validEnvs := []string{"development", "testing", "production"}
	validEnv := false
	for _, env := range validEnvs {
		if config.App.Environment == env {
			validEnv = true
			break
		}
	}
	if !validEnv {
		errs = append(errs, fmt.Sprintf("invalid environment: %s", config.App.Environment))
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}
