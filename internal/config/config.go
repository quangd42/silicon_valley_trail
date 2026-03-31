package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	WeatherAPIKey     string
	WeatherAPIUseMock bool
	WeatherAPITimeout time.Duration
	SavePath          string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}
	cfg := Config{
		WeatherAPIKey:     os.Getenv("WEATHERAPI_KEY"),
		WeatherAPIUseMock: os.Getenv("WEATHERAPI_KEY") == "",
		WeatherAPITimeout: 3 * time.Second,
		SavePath:          os.Getenv("SAVE_PATH"),
	}

	if raw, ok := os.LookupEnv("WEATHERAPI_MOCK"); ok && raw != "" {
		useMock, err := strconv.ParseBool(raw)
		if err != nil {
			return Config{}, fmt.Errorf("parse WEATHERAPI_MOCK: %w", err)
		}
		cfg.WeatherAPIUseMock = useMock
	}

	if raw, ok := os.LookupEnv("WEATHERAPI_TIMEOUT_MS"); ok && raw != "" {
		timeout, err := strconv.Atoi(raw)
		if err != nil {
			return Config{}, fmt.Errorf("parse WEATHERAPI_TIMEOUT_MS: %w", err)
		}
		cfg.WeatherAPITimeout = time.Duration(timeout) * time.Millisecond
	}

	return cfg, nil
}
