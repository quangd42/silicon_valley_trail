package weather

import (
	"context"

	"github.com/quangd42/silicon_valley_trail/internal/model"
)

type MockService struct {
	data  []model.WeatherKind
	index int
}

func NewMockService(data []model.WeatherKind) *MockService {
	return &MockService{
		data:  data,
		index: 0,
	}
}

func DefaultMockService() *MockService {
	return NewMockService([]model.WeatherKind{
		model.WeatherClear,
		model.WeatherRainy,
		model.WeatherFog,
		model.WeatherCloudy,
	})
}

func (s *MockService) WeatherAt(_ context.Context, _ model.Location) (model.WeatherKind, error) {
	out := s.data[s.index]
	s.index++
	if s.index == len(s.data) {
		s.index = 0
	}
	return out, nil
}
