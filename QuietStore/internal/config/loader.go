package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

func Load() (*Config, error) {
	cfg := &Config{}

	if err := loadStruct(&cfg.Server, ""); err != nil {
		return nil, fmt.Errorf("failed to load server config: %v", err)
	}
	if err := loadStruct(&cfg.App, ""); err != nil {
		return nil, fmt.Errorf("loading app config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	return cfg, nil
}

func loadStruct(anyConfig interface{}, prefix string) error {
	v := reflect.ValueOf(anyConfig).Elem() // represents dereferenced address of config struct (so the actual struct)
	t := v.Type()                          // represents type of config struct i.e. Config, ServerConfig, etc.

	for i := 0; i < v.NumField(); i++ { // for each field in the struct
		field := t.Field(i) // represents field type and name i.e. ReadTimeout time.Duration
		value := v.Field(i)

		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		envVar := envTag
		if prefix != "" {
			envVar = prefix + "_" + envTag
		}

		envValue := os.Getenv(envVar)

		if envValue == "" {
			defaultTag := field.Tag.Get("default")
			if defaultTag != "" {
				envValue = defaultTag
			} else if field.Tag.Get("required") == "true" {
				return fmt.Errorf("required environment variable %s is not set", envVar)
			}
		}

		if err := setField(value, envValue); err != nil {
			return fmt.Errorf("failed to set field %s: %w", field.Name, err)
		}
	}

	return nil
}

func setField(field reflect.Value, value string) error {
	if value == "" && field.Kind() != reflect.String {
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Int, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intValue)

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolValue)

	case reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatValue)

	default:
		if field.Type() == reflect.TypeOf(time.Duration(0)) {
			duration, err := time.ParseDuration(value)
			if err != nil {
				return fmt.Errorf("invalid duration value %s: %w", value, err)
			}
			field.Set(reflect.ValueOf(duration))
			return nil
		}
		return fmt.Errorf("unsupported field type %s for value %s", field.Type(), value)
	}

	return nil
}
