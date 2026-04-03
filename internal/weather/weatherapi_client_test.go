package weather

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/quangd42/silicon_valley_trail/internal/model"
)

func TestWeatherAPIService_Current(t *testing.T) {
	t.Run("returns weather and caches by location", func(t *testing.T) {
		var requestCount int
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			if r.URL.Path != "/v1/current.json" {
				t.Fatalf("path = %q, want %q", r.URL.Path, "/v1/current.json")
			}
			if got := r.URL.Query().Get("aqi"); got != "no" {
				t.Fatalf("aqi = %q, want %q", got, "no")
			}
			if got := r.URL.Query().Get("key"); got != "test-key" {
				t.Fatalf("key = %q, want %q", got, "test-key")
			}
			if got := r.URL.Query().Get("q"); got != "37.368900,-122.035300" {
				t.Fatalf("q = %q, want %q", got, "37.368900,-122.035300")
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"current":{"condition":{"code":1035}}}`))
		}))
		defer server.Close()

		service := NewWeatherAPIService("test-key", server.URL+"/v1/current.json?aqi=no", 0)
		loc := model.Location{
			ID:        "sunnyvale",
			Latitude:  37.3689,
			Longitude: -122.0353,
		}

		first, err := service.WeatherAt(context.Background(), loc)
		if err != nil {
			t.Fatalf("first Current() error = %v, want nil", err)
		}
		second, err := service.WeatherAt(context.Background(), loc)
		if err != nil {
			t.Fatalf("second Current() error = %v, want nil", err)
		}

		if first != model.WeatherRainy {
			t.Fatalf("first Current() = %v, want %v", first, model.WeatherRainy)
		}
		if second != model.WeatherRainy {
			t.Fatalf("second Current() = %v, want %v", second, model.WeatherRainy)
		}
		if requestCount != 1 {
			t.Fatalf("requestCount = %d, want 1", requestCount)
		}
		if cached, ok := service.cache[loc.ID]; !ok || cached != model.WeatherRainy {
			t.Fatalf("cache[%q] = %v, %v; want %v, true", loc.ID, cached, ok, model.WeatherRainy)
		}
	})

	t.Run("returns missing location error", func(t *testing.T) {
		service := NewWeatherAPIService("test-key", "http://example.invalid", 0)

		got, err := service.WeatherAt(context.Background(), model.Location{})
		if got != model.WeatherUnknown {
			t.Fatalf("Current() weather = %v, want %v", got, model.WeatherUnknown)
		}
		if !errors.Is(err, ErrMissingLocation) {
			t.Fatalf("Current() error = %v, want %v", err, ErrMissingLocation)
		}
	})

	t.Run("returns request failed on non-2xx response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "bad gateway", http.StatusBadGateway)
		}))
		defer server.Close()

		service := NewWeatherAPIService("test-key", server.URL+"/v1/current.json?aqi=no", 0)
		got, err := service.WeatherAt(context.Background(), model.Location{
			ID:        "san-jose",
			Latitude:  37.3394,
			Longitude: -121.8939,
		})
		if got != model.WeatherUnknown {
			t.Fatalf("Current() weather = %v, want %v", got, model.WeatherUnknown)
		}
		if !errors.Is(err, ErrRequestFailed) {
			t.Fatalf("Current() error = %v, want %v", err, ErrRequestFailed)
		}
	})

	t.Run("returns invalid response shape on malformed json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{not-json`))
		}))
		defer server.Close()

		service := NewWeatherAPIService("test-key", server.URL+"/v1/current.json?aqi=no", 0)
		got, err := service.WeatherAt(context.Background(), model.Location{
			ID:        "san-jose",
			Latitude:  37.3394,
			Longitude: -121.8939,
		})
		if got != model.WeatherUnknown {
			t.Fatalf("Current() weather = %v, want %v", got, model.WeatherUnknown)
		}
		if !errors.Is(err, ErrInvalidResponseShape) {
			t.Fatalf("Current() error = %v, want %v", err, ErrInvalidResponseShape)
		}
	})

	t.Run("returns timeout error when request exceeds client timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{"current":{"condition":{"code":1000}}}`)
		}))
		defer server.Close()

		service := NewWeatherAPIService("test-key", server.URL+"/v1/current.json?aqi=no", 10*time.Millisecond)
		got, err := service.WeatherAt(context.Background(), model.Location{
			ID:        "san-jose",
			Latitude:  37.3394,
			Longitude: -121.8939,
		})
		if got != model.WeatherUnknown {
			t.Fatalf("Current() weather = %v, want %v", got, model.WeatherUnknown)
		}
		if err == nil {
			t.Fatal("Current() error = nil, want timeout error")
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Current() error = %v, want deadline exceeded", err)
		}
	})
}

func TestToWeatherKind(t *testing.T) {
	tests := []struct {
		code int
		want model.WeatherKind
	}{
		{code: 1000, want: model.WeatherClear},
		{code: 1006, want: model.WeatherCloudy},
		{code: 1135, want: model.WeatherFog},
		{code: 1183, want: model.WeatherRainy},
	}

	for _, tt := range tests {
		if got := toWeatherKind(tt.code); got != tt.want {
			t.Fatalf("toWeatherKind(%d) = %q, want %q", tt.code, got, tt.want)
		}
	}
}
