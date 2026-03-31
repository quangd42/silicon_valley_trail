package weather

import (
	"context"
	"errors"
	"testing"

	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type stubService struct {
	kind model.WeatherKind
	err  error
}

func (s stubService) Current(context.Context, model.Location) (model.WeatherKind, error) {
	return s.kind, s.err
}

func TestWeatherService(t *testing.T) {
	t.Run("uses mock service when api key is missing", func(t *testing.T) {
		svc := NewWeatherService("", false, 0)

		got, err := svc.WeatherAt(context.Background(), model.Location{ID: "san-jose"})
		if err != nil {
			t.Fatalf("Current() error = %v", err)
		}
		if got != model.WeatherClear {
			t.Fatalf("Current() = %q, want %q", got, model.WeatherClear)
		}
	})

	t.Run("falls back to mock on remote error", func(t *testing.T) {
		svc := &WeatherService{
			mock:   stubService{kind: model.WeatherFog}.Current,
			remote: stubService{err: errors.New("boom")}.Current,
		}

		got, err := svc.WeatherAt(context.Background(), model.Location{ID: "san-jose"})
		if err != nil {
			t.Fatalf("Current() error = %v, want nil", err)
		}
		if got != model.WeatherFog {
			t.Fatalf("Current() = %q, want %q", got, model.WeatherFog)
		}
	})
}
