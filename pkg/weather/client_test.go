package weather

import (
	"errors"
	"strings"
	"testing"

	"github.com/JSGette/agent_summit_bazel_workshop/internal/testutils"
	"github.com/JSGette/agent_summit_bazel_workshop/pkg/models"
)

func TestClient_GetWeatherByCoordinates(t *testing.T) {
	tests := []struct {
		name           string
		lat            float64
		lon            float64
		city           string
		country        string
		mockResponse   string
		mockStatusCode int
		mockError      error
		wantError      bool
		wantTemp       float64
		wantCondition  models.WeatherCondition
	}{
		{
			name:           "successful weather request",
			lat:            48.7758,
			lon:            9.1829,
			city:           "Stuttgart",
			country:        "Germany",
			mockResponse:   testutils.OpenMeteoWeatherResponse,
			mockStatusCode: 200,
			wantError:      false,
			wantTemp:       22.5,
			wantCondition:  models.Cloudy,
		},
		{
			name:           "API returns 500 error",
			lat:            48.7758,
			lon:            9.1829,
			city:           "Stuttgart",
			country:        "Germany",
			mockResponse:   testutils.APIErrorResponse,
			mockStatusCode: 500,
			wantError:      true,
		},
		{
			name:      "network error",
			lat:       48.7758,
			lon:       9.1829,
			city:      "Stuttgart",
			country:   "Germany",
			mockError: errors.New("network error"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := testutils.NewMockHTTPClient()
			client := NewClient(mockClient)

			// Prepare expected URL
			expectedURL := "https://api.open-meteo.com/v1/forecast?current=temperature_2m%2Cweather_code%2Cis_day&latitude=48.7758&longitude=9.1829&timezone=auto"

			if tt.mockError != nil {
				mockClient.AddError(expectedURL, tt.mockError)
			} else {
				mockClient.AddResponse(expectedURL, tt.mockStatusCode, tt.mockResponse)
			}

			result, err := client.GetWeatherByCoordinates(tt.lat, tt.lon, tt.city, tt.country)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("Expected result, but got nil")
				return
			}

			if result.Temperature != tt.wantTemp {
				t.Errorf("Expected temperature %v, got %v", tt.wantTemp, result.Temperature)
			}

			if result.Condition != tt.wantCondition {
				t.Errorf("Expected condition %v, got %v", tt.wantCondition, result.Condition)
			}

			if result.City != tt.city {
				t.Errorf("Expected city %v, got %v", tt.city, result.City)
			}

			if result.Country != tt.country {
				t.Errorf("Expected country %v, got %v", tt.country, result.Country)
			}
		})
	}
}

func TestClient_GetWeatherByCity(t *testing.T) {
	tests := []struct {
		name              string
		city              string
		mockGeocodeResp   string
		mockWeatherResp   string
		mockGeocodeStatus int
		mockWeatherStatus int
		mockGeocodeError  error
		mockWeatherError  error
		wantError         bool
		wantCity          string
	}{
		{
			name:              "successful request for Stuttgart",
			city:              "Stuttgart",
			mockGeocodeResp:   testutils.OpenMeteoGeocodeResponse,
			mockWeatherResp:   testutils.OpenMeteoWeatherResponse,
			mockGeocodeStatus: 200,
			mockWeatherStatus: 200,
			wantError:         false,
			wantCity:          "Stuttgart",
		},
		{
			name:              "city not found",
			city:              "NonexistentCity",
			mockGeocodeResp:   testutils.OpenMeteoGeocodeNotFound,
			mockGeocodeStatus: 200,
			wantError:         true,
		},
		{
			name:             "geocoding API error",
			city:             "Stuttgart",
			mockGeocodeError: errors.New("geocoding error"),
			wantError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := testutils.NewMockHTTPClient()
			client := NewClient(mockClient)

			// Setup geocoding mock
			geocodeURL := "https://geocoding-api.open-meteo.com/v1/search?count=1&format=json&language=en&name=" + tt.city
			if tt.mockGeocodeError != nil {
				mockClient.AddError(geocodeURL, tt.mockGeocodeError)
			} else {
				mockClient.AddResponse(geocodeURL, tt.mockGeocodeStatus, tt.mockGeocodeResp)
			}

			// Setup weather mock if geocoding succeeds
			if !tt.wantError && tt.mockGeocodeError == nil && tt.mockGeocodeStatus == 200 {
				weatherURL := "https://api.open-meteo.com/v1/forecast?current=temperature_2m%2Cweather_code%2Cis_day&latitude=48.7758&longitude=9.1829&timezone=auto"
				if tt.mockWeatherError != nil {
					mockClient.AddError(weatherURL, tt.mockWeatherError)
				} else {
					mockClient.AddResponse(weatherURL, tt.mockWeatherStatus, tt.mockWeatherResp)
				}
			}

			result, err := client.GetWeatherByCity(tt.city)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("Expected result, but got nil")
				return
			}

			if result.City != tt.wantCity {
				t.Errorf("Expected city %v, got %v", tt.wantCity, result.City)
			}
		})
	}
}

func TestClient_GetWeather(t *testing.T) {
	tests := []struct {
		name      string
		location  string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "empty location",
			location:  "",
			wantError: true,
			errorMsg:  "Location cannot be empty",
		},
		{
			name:      "valid location",
			location:  "Stuttgart",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := testutils.NewMockHTTPClient()
			client := NewClient(mockClient)

			// For valid locations, mock the responses
			if !tt.wantError {
				geocodeURL := "https://geocoding-api.open-meteo.com/v1/search?count=1&format=json&language=en&name=" + tt.location
				mockClient.AddResponse(geocodeURL, 200, testutils.OpenMeteoGeocodeResponse)

				weatherURL := "https://api.open-meteo.com/v1/forecast?current=temperature_2m%2Cweather_code%2Cis_day&latitude=48.7758&longitude=9.1829&timezone=auto"
				mockClient.AddResponse(weatherURL, 200, testutils.OpenMeteoWeatherResponse)
			}

			_, err := client.GetWeather(tt.location)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error, but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', got: %v", tt.errorMsg, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
