package config

import (
	"errors"
	"fmt"
	"strings"
)

func (config *Config) Validate() error {
	var errs []string

	if config.Server.Port < 1 || config.Server.Port > 65535 {
		errs = append(errs, "server port must be between 1 and 65535")
	}

	validEnvs := []string{"development", "testing", "production"}
	if !contains(validEnvs, config.App.Environment) {
		errs = append(errs, fmt.Sprintf("invalid environment: %s", config.App.Environment))
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

func contains(list []string, val string) bool {
	for _, item := range list {
		if item == val {
			return true
		}
	}
	return false
}
