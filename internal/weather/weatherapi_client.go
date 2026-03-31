package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/quangd42/silicon_valley_trail/internal/model"
)

const WeatherAPIEndpoint = "https://api.weatherapi.com/v1/current.json?aqi=no"

var (
	ErrMissingLocation      = errors.New("location ID missing")
	ErrRequestFailed        = errors.New("request failed")
	ErrInvalidResponseShape = errors.New("response shape mismatched")
)

// WeatherAPIService requests live weather data from defined remote and caches
// it for the duration of the game session. No other cache invalidation scheme
// is defined.
type WeatherAPIService struct {
	apiKey  string
	baseURL string
	cache   map[string]model.WeatherKind
	client  *http.Client
}

func NewWeatherAPIService(apiKey, baseURL string, timeout time.Duration) *WeatherAPIService {
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	if baseURL == "" {
		baseURL = WeatherAPIEndpoint
	}
	return &WeatherAPIService{
		apiKey:  apiKey,
		baseURL: baseURL,
		cache:   make(map[string]model.WeatherKind),
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *WeatherAPIService) Current(ctx context.Context, loc model.Location) (model.WeatherKind, error) {
	if cached, ok := s.cache[loc.ID]; ok {
		return cached, nil
	}
	if loc.ID == "" {
		return model.WeatherUnknown, ErrMissingLocation
	}
	endpoint := s.buildURL(loc)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return model.WeatherUnknown, err
	}
	res, err := s.client.Do(req)
	if err != nil {
		return model.WeatherUnknown, err
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return model.WeatherUnknown, ErrRequestFailed
	}
	type Response struct {
		Current struct {
			Condition struct {
				Code int `json:"code"`
			} `json:"condition"`
		} `json:"current"`
	}
	var data Response
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return model.WeatherUnknown, ErrInvalidResponseShape
	}
	weatherKind := toWeatherKind(data.Current.Condition.Code)
	s.cache[loc.ID] = weatherKind
	return weatherKind, nil
}

func (s *WeatherAPIService) buildURL(loc model.Location) string {
	u, err := url.Parse(s.baseURL)
	if err != nil {
		panic("invalid weatherapi endpoint")
	}
	q := u.Query()
	q.Set("key", s.apiKey)
	q.Set("q", fmt.Sprintf("%f,%f", loc.Latitude, loc.Longitude))
	u.RawQuery = q.Encode()
	return u.String()
}

// For simplicity and the fact that most codes points to a rainy weather condition
// the catch all case returns `WeatherRainy`.
// https://www.weatherapi.com/docs/#weather-icons
func toWeatherKind(code int) model.WeatherKind {
	switch {
	case code == 1000:
		return model.WeatherClear
	case code >= 1003 && code <= 1009:
		return model.WeatherCloudy
	case code == 1135 || code == 1147:
		return model.WeatherFog
	default:
		return model.WeatherRainy
	}
}
