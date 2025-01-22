package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	LogLevel string
	DataPath string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found (or unable to load). Continuing with system environment variables. Error: %v", err)
	}

	cfg := &Config{
		Port:     getEnvWithDefault("PORT", "8080"),
		LogLevel: getEnvWithDefault("LOG_LEVEL", "info"),
		DataPath: getEnvWithDefault("DATA_PATH", "data/input.txt"),
	}

	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func validate(cfg *Config) error {
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"error": true,
	}

	if !validLogLevels[cfg.LogLevel] {
		return fmt.Errorf("invalid log level: %s", cfg.LogLevel)
	}

	return nil
}
