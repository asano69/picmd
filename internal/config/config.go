// Package config loads the configuration for myapp serve from
// environment variables.
package config

import (
	"os"
	"strconv"

	"github.com/asano69/myapp/internal/errs"
)

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host string
	Port int
}

// DataConfig holds data storage settings.
type DataConfig struct {
	Root string
}
type Config struct {
	Server ServerConfig
	Data   DataConfig
}

// Load reads configuration from environment variables, applying defaults
// for any variable that is unset.
//
// Recognised variables:
//
//	SERVER_HOST         default "0.0.0.0"
//	SERVER_PORT         default 3000
//	DATA_ROOT           default "."

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host: envString("SERVER_HOST", "0.0.0.0"),
			Port: 3000,
		},
		Data: DataConfig{
			Root: envString("DATA_ROOT", "."),
		},
	}

	port, err := envInt("SERVER_PORT", cfg.Server.Port)
	if err != nil {
		return nil, err
	}
	cfg.Server.Port = port

	return cfg, nil
}

// envString returns the value of the environment variable key, or fallback
// if it is unset or empty.
func envString(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

// envInt returns the integer value of the environment variable key, or
// fallback if it is unset.
func envInt(key string, fallback int) (int, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, errs.Newf("invalid %s: %v", key, err)
	}
	return n, nil
}

// envFloat returns the float64 value of the environment variable key, or
// fallback if it is unset.
func envFloat(key string, fallback float64) (float64, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback, nil
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, errs.Newf("invalid %s: %v", key, err)
	}
	return f, nil
}
