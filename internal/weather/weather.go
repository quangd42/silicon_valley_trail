package weather

import (
	"context"
	"time"

	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type weatherProviderFn func(context.Context, model.Location) (model.WeatherKind, error)

type WeatherService struct {
	mock     weatherProviderFn
	remote   weatherProviderFn
	mockOnly bool
}

func NewWeatherService(api string, mock bool, timeout time.Duration) *WeatherService {
	mockSvc := NewMockService()
	if api == "" || mock {
		return &WeatherService{
			mock:     mockSvc.Current,
			mockOnly: true,
		}
	}
	remoteSvc := NewWeatherAPIService(api, "", timeout)
	return &WeatherService{
		mockOnly: false,
		mock:     mockSvc.Current,
		remote:   remoteSvc.Current,
	}
}

func (s *WeatherService) WeatherAt(ctx context.Context, l model.Location) (model.WeatherKind, error) {
	if s.mockOnly {
		return s.mock(ctx, l)
	}
	out, err := s.remote(ctx, l)
	if err != nil {
		return s.mock(ctx, l)
	}
	return out, nil
}
