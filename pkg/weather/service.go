package weather

import (
	"fmt"
	"log"
	"time"

	"github.com/JSGette/agent_summit_bazel_workshop/pkg/models"
)

// Service provides high-level weather operations with caching and logging
type Service struct {
	client *Client
}

// NewService creates a new weather service
func NewService(httpClient HTTPClient) *Service {
	return &Service{
		client: NewClient(httpClient),
	}
}

// GetCurrentWeather fetches current weather for a location with enhanced error handling
func (s *Service) GetCurrentWeather(location string) (*models.WeatherResponse, error) {
	start := time.Now()

	log.Printf("Fetching weather for location: %s", location)

	weather, err := s.client.GetWeather(location)
	if err != nil {
		log.Printf("Error fetching weather for %s: %v", location, err)
		return nil, err
	}

	duration := time.Since(start)
	log.Printf("Successfully fetched weather for %s in %v", location, duration)

	return weather, nil
}

// GetWeatherSummary returns a human-readable weather summary
func (s *Service) GetWeatherSummary(location string) (string, error) {
	weather, err := s.GetCurrentWeather(location)
	if err != nil {
		return "", err
	}

	timeOfDay := "during the day"
	if !weather.IsDay {
		timeOfDay = "during the night"
	}

	summary := fmt.Sprintf(
		"Current weather in %s, %s: %.1f°C, %s %s. Last updated: %s",
		weather.City,
		weather.Country,
		weather.Temperature,
		weather.Description,
		timeOfDay,
		weather.Metadata.Timestamp.Format("15:04 MST"),
	)

	return summary, nil
}

// ValidateLocation checks if a location string is valid
func (s *Service) ValidateLocation(location string) error {
	if location == "" {
		return models.NewAPIError("Weather Service", "Location cannot be empty", 400)
	}

	if len(location) < 2 {
		return models.NewAPIError("Weather Service", "Location must be at least 2 characters long", 400)
	}

	if len(location) > 100 {
		return models.NewAPIError("Weather Service", "Location must be less than 100 characters", 400)
	}

	return nil
}

// GetWeatherWithValidation fetches weather with input validation
func (s *Service) GetWeatherWithValidation(location string) (*models.WeatherResponse, error) {
	if err := s.ValidateLocation(location); err != nil {
		return nil, err
	}

	return s.GetCurrentWeather(location)
}
