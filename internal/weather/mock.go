package weather

import (
	"context"

	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type MockService struct {
	data  []model.WeatherKind
	index int
}

func NewMockService() *MockService {
	return &MockService{
		data: []model.WeatherKind{
			model.WeatherClear,
			model.WeatherRainy,
			model.WeatherFog,
			model.WeatherCloudy,
		},
		index: 0,
	}
}

func (s *MockService) WeatherAt(_ context.Context, _ model.Location) (model.WeatherKind, error) {
	out := s.data[s.index]
	s.index++
	if s.index == len(s.data) {
		s.index = 0
	}
	return out, nil
}
