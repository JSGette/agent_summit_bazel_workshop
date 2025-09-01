package weather

import (
	"strings"
	"testing"

	"github.com/JSGette/agent_summit_bazel_workshop/internal/testutils"
)

func TestService_GetCurrentWeather(t *testing.T) {
	tests := []struct {
		name           string
		location       string
		mockResponse   string
		mockStatusCode int
		wantError      bool
	}{
		{
			name:           "successful weather request",
			location:       "Stuttgart",
			mockResponse:   testutils.OpenMeteoWeatherResponse,
			mockStatusCode: 200,
			wantError:      false,
		},
		{
			name:      "empty location",
			location:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := testutils.NewMockHTTPClient()
			service := NewService(mockClient)

			// Mock geocoding response for valid locations
			if tt.location != "" {
				geocodeURL := "https://geocoding-api.open-meteo.com/v1/search?count=1&format=json&language=en&name=" + tt.location
				mockClient.AddResponse(geocodeURL, 200, testutils.OpenMeteoGeocodeResponse)

				weatherURL := "https://api.open-meteo.com/v1/forecast?current=temperature_2m%2Cweather_code%2Cis_day&latitude=48.7758&longitude=9.1829&timezone=auto"
				mockClient.AddResponse(weatherURL, tt.mockStatusCode, tt.mockResponse)
			}

			result, err := service.GetCurrentWeather(tt.location)

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
			}
		})
	}
}

func TestService_GetWeatherSummary(t *testing.T) {
	mockClient := testutils.NewMockHTTPClient()
	service := NewService(mockClient)

	// Mock successful response
	geocodeURL := "https://geocoding-api.open-meteo.com/v1/search?count=1&format=json&language=en&name=Stuttgart"
	mockClient.AddResponse(geocodeURL, 200, testutils.OpenMeteoGeocodeResponse)

	weatherURL := "https://api.open-meteo.com/v1/forecast?current=temperature_2m%2Cweather_code%2Cis_day&latitude=48.7758&longitude=9.1829&timezone=auto"
	mockClient.AddResponse(weatherURL, 200, testutils.OpenMeteoWeatherResponse)

	summary, err := service.GetWeatherSummary("Stuttgart")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expectedParts := []string{"Stuttgart", "Germany", "22.5°C", "Overcast"}
	for _, part := range expectedParts {
		if !strings.Contains(summary, part) {
			t.Errorf("Expected summary to contain '%s', got: %s", part, summary)
		}
	}
}

func TestService_ValidateLocation(t *testing.T) {
	tests := []struct {
		name      string
		location  string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid location",
			location:  "Stuttgart",
			wantError: false,
		},
		{
			name:      "empty location",
			location:  "",
			wantError: true,
			errorMsg:  "Location cannot be empty",
		},
		{
			name:      "too short location",
			location:  "S",
			wantError: true,
			errorMsg:  "at least 2 characters",
		},
		{
			name:      "too long location",
			location:  strings.Repeat("a", 101),
			wantError: true,
			errorMsg:  "less than 100 characters",
		},
		{
			name:      "minimum valid length",
			location:  "NY",
			wantError: false,
		},
		{
			name:      "maximum valid length",
			location:  strings.Repeat("a", 100),
			wantError: false,
		},
	}

	service := NewService(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateLocation(tt.location)

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

func TestService_GetWeatherWithValidation(t *testing.T) {
	tests := []struct {
		name      string
		location  string
		wantError bool
	}{
		{
			name:      "valid location",
			location:  "Stuttgart",
			wantError: false,
		},
		{
			name:      "invalid location - empty",
			location:  "",
			wantError: true,
		},
		{
			name:      "invalid location - too short",
			location:  "S",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := testutils.NewMockHTTPClient()
			service := NewService(mockClient)

			// Mock successful API responses for valid locations
			if !tt.wantError {
				geocodeURL := "https://geocoding-api.open-meteo.com/v1/search?count=1&format=json&language=en&name=" + tt.location
				mockClient.AddResponse(geocodeURL, 200, testutils.OpenMeteoGeocodeResponse)

				weatherURL := "https://api.open-meteo.com/v1/forecast?current=temperature_2m%2Cweather_code%2Cis_day&latitude=48.7758&longitude=9.1829&timezone=auto"
				mockClient.AddResponse(weatherURL, 200, testutils.OpenMeteoWeatherResponse)
			}

			_, err := service.GetWeatherWithValidation(tt.location)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestNewService(t *testing.T) {
	t.Run("with custom client", func(t *testing.T) {
		mockClient := testutils.NewMockHTTPClient()
		service := NewService(mockClient)
		if service == nil {
			t.Errorf("Expected service, but got nil")
		}
		if service.client == nil {
			t.Errorf("Expected client to be set")
		}
	})

	t.Run("with nil client", func(t *testing.T) {
		service := NewService(nil)
		if service == nil {
			t.Errorf("Expected service, but got nil")
		}
		if service.client == nil {
			t.Errorf("Expected default client to be set")
		}
	})
}
