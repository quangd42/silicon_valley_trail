package weather

import (
	"context"
	"time"

	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type Service interface {
	Current(context.Context, model.Location) (model.WeatherKind, error)
}

type WeatherService struct {
	mock     Service
	remote   Service
	mockOnly bool
}

func NewWeatherService(api string, mock bool, timeout time.Duration) Service {
	if api == "" || mock {
		return &WeatherService{
			mock:     NewMockService(),
			mockOnly: true,
		}
	}
	return &WeatherService{
		mockOnly: false,
		mock:     NewMockService(),
		remote:   NewWeatherAPIService(api, "", timeout),
	}
}

func (s *WeatherService) Current(ctx context.Context, l model.Location) (model.WeatherKind, error) {
	if s.mockOnly {
		return s.mock.Current(ctx, l)
	}
	out, err := s.remote.Current(ctx, l)
	if err != nil {
		return s.mock.Current(ctx, l)
	}
	return out, nil
}
