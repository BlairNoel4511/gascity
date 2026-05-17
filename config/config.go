// Package config provides configuration management for gascity.
// It supports loading configuration from environment variables and config files.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration values.
type Config struct {
	// Server configuration
	Server ServerConfig

	// Database configuration
	Database DatabaseConfig

	// Application configuration
	App AppConfig
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	ConnLifetime time.Duration
}

// AppConfig holds general application configuration.
type AppConfig struct {
	Environment string
	LogLevel    string
	Debug       bool
}

// Load reads configuration from environment variables and returns a Config.
// Environment variables take precedence over default values.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			// Increased timeouts from 15s to 30s — upstream defaults felt too tight
			// for slower queries I've been testing locally.
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			DSN:          getEnv("DATABASE_DSN", ""),
			MaxOpenConns: getEnvAsInt("DATABASE_MAX_OPEN_CONNS", 25),
			// Bumped MaxIdleConns from 5 to 10 — noticed connection churn in local
			// load tests; keeping more idle connections ready helps throughput.
			MaxIdleConns: getEnvAsInt("DATABASE_MAX_IDLE_CONNS", 10),
			ConnLifetime: getEnvAsDuration("DATABASE_CONN_LIFETIME", 5*time.Minute),
		},
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
			Debug:       getEnvAsBool("APP_DEBUG", false),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// validate checks that required configuration values are set and valid.
func (c *Config) validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("server port %d is out of valid range (1-65535)", c.Server.Port)
	}
	return nil
}

// Addr returns the full server address string.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

// getEnvAsInt retrieves an environment variable as an integer or returns a default value.
func getEnvAsInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); e