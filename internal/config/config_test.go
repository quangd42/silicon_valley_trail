package config

import (
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	t.Run("parses explicit weather config", func(t *testing.T) {
		t.Setenv("WEATHERAPI_KEY", "test-key")
		t.Setenv("WEATHERAPI_MOCK", "false")
		t.Setenv("WEATHERAPI_TIMEOUT_MS", "4500")
		t.Setenv("SAVE_PATH", "custom-save.json")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if cfg.WeatherAPIKey != "test-key" {
			t.Fatalf("WeatherAPIKey = %q, want %q", cfg.WeatherAPIKey, "test-key")
		}
		if cfg.WeatherAPIUseMock {
			t.Fatalf("WeatherAPIUseMock = %v, want false", cfg.WeatherAPIUseMock)
		}
		if cfg.WeatherAPITimeout != 4500*time.Millisecond {
			t.Fatalf("WeatherAPITimeout = %v, want %v", cfg.WeatherAPITimeout, 4500*time.Millisecond)
		}
		if cfg.SavePath != "custom-save.json" {
			t.Fatalf("SavePath = %q, want %q", cfg.SavePath, "custom-save.json")
		}
	})

	t.Run("defaults to mock mode when weather env is missing", func(t *testing.T) {
		t.Setenv("WEATHERAPI_KEY", "")
		t.Setenv("WEATHERAPI_MOCK", "")
		t.Setenv("WEATHERAPI_TIMEOUT_MS", "")
		t.Setenv("SAVE_PATH", "")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if !cfg.WeatherAPIUseMock {
			t.Fatalf("WeatherAPIUseMock = %v, want true", cfg.WeatherAPIUseMock)
		}
		if cfg.WeatherAPITimeout != 3*time.Second {
			t.Fatalf("WeatherAPITimeout = %v, want %v", cfg.WeatherAPITimeout, 3*time.Second)
		}
		if cfg.SavePath != "" {
			t.Fatalf("SavePath = %q, want empty", cfg.SavePath)
		}
	})
}
