package weather

import (
	"errors"
	"strings"
	"testing"

	"github.com/JSGette/agent_summit_bazel_workshop/internal/testutils"
)

func TestGeocoder_GetCoordinates(t *testing.T) {
	tests := []struct {
		name           string
		city           string
		mockResponse   string
		mockStatusCode int
		mockError      error
		wantError      bool
		wantLat        float64
		wantLon        float64
		wantCountry    string
	}{
		{
			name:           "successful geocoding",
			city:           "Stuttgart",
			mockResponse:   testutils.OpenMeteoGeocodeResponse,
			mockStatusCode: 200,
			wantError:      false,
			wantLat:        48.7758,
			wantLon:        9.1829,
			wantCountry:    "Germany",
		},
		{
			name:           "city not found",
			city:           "NonexistentCity",
			mockResponse:   testutils.OpenMeteoGeocodeNotFound,
			mockStatusCode: 200,
			wantError:      true,
		},
		{
			name:      "empty city name",
			city:      "",
			wantError: true,
		},
		{
			name:      "whitespace only city name",
			city:      "   ",
			wantError: true,
		},
		{
			name:           "API returns 500 error",
			city:           "Stuttgart",
			mockResponse:   testutils.APIErrorResponse,
			mockStatusCode: 500,
			wantError:      true,
		},
		{
			name:      "network error",
			city:      "Stuttgart",
			mockError: errors.New("network error"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := testutils.NewMockHTTPClient()
			geocoder := NewGeocoder(mockClient)

			if tt.city != "" && strings.TrimSpace(tt.city) != "" {
				expectedURL := "https://geocoding-api.open-meteo.com/v1/search?count=1&format=json&language=en&name=" + tt.city

				if tt.mockError != nil {
					mockClient.AddError(expectedURL, tt.mockError)
				} else {
					mockClient.AddResponse(expectedURL, tt.mockStatusCode, tt.mockResponse)
				}
			}

			coords, country, err := geocoder.GetCoordinates(tt.city)

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

			if coords == nil {
				t.Errorf("Expected coordinates, but got nil")
				return
			}

			if coords.Latitude != tt.wantLat {
				t.Errorf("Expected latitude %v, got %v", tt.wantLat, coords.Latitude)
			}

			if coords.Longitude != tt.wantLon {
				t.Errorf("Expected longitude %v, got %v", tt.wantLon, coords.Longitude)
			}

			if country != tt.wantCountry {
				t.Errorf("Expected country %v, got %v", tt.wantCountry, country)
			}
		})
	}
}

func TestGeocoder_GetCoordinatesWithCache(t *testing.T) {
	tests := []struct {
		name        string
		city        string
		expectCache bool
		wantLat     float64
		wantLon     float64
		wantCountry string
	}{
		{
			name:        "cached city - stuttgart",
			city:        "Stuttgart",
			expectCache: true,
			wantLat:     48.7758,
			wantLon:     9.1829,
			wantCountry: "Germany",
		},
		{
			name:        "cached city - case insensitive",
			city:        "BERLIN",
			expectCache: true,
			wantLat:     52.5200,
			wantLon:     13.4050,
			wantCountry: "Germany",
		},
		{
			name:        "cached city - with spaces",
			city:        "  New York  ",
			expectCache: true,
			wantLat:     40.7128,
			wantLon:     -74.0060,
			wantCountry: "United States",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := testutils.NewMockHTTPClient()
			geocoder := NewGeocoder(mockClient)

			// If not expecting cache, mock the API response
			if !tt.expectCache {
				expectedURL := "https://geocoding-api.open-meteo.com/v1/search?count=1&format=json&language=en&name=" + tt.city
				mockClient.AddResponse(expectedURL, 200, testutils.OpenMeteoGeocodeResponse)
			}

			coords, country, err := geocoder.GetCoordinatesWithCache(tt.city)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if coords == nil {
				t.Errorf("Expected coordinates, but got nil")
				return
			}

			if coords.Latitude != tt.wantLat {
				t.Errorf("Expected latitude %v, got %v", tt.wantLat, coords.Latitude)
			}

			if coords.Longitude != tt.wantLon {
				t.Errorf("Expected longitude %v, got %v", tt.wantLon, coords.Longitude)
			}

			if country != tt.wantCountry {
				t.Errorf("Expected country %v, got %v", tt.wantCountry, country)
			}

			// Verify that API was not called for cached cities
			if tt.expectCache {
				apiCallCount := 0
				for url, count := range mockClient.CallCount {
					if strings.Contains(url, "geocoding-api.open-meteo.com") {
						apiCallCount += count
					}
				}
				if apiCallCount > 0 {
					t.Errorf("Expected cache hit, but API was called %d times", apiCallCount)
				}
			}
		})
	}
}

func TestNewGeocoder(t *testing.T) {
	t.Run("with nil client", func(t *testing.T) {
		geocoder := NewGeocoder(nil)
		if geocoder == nil {
			t.Errorf("Expected geocoder, but got nil")
		}
		if geocoder.client == nil {
			t.Errorf("Expected default client to be set")
		}
	})

	t.Run("with custom client", func(t *testing.T) {
		mockClient := testutils.NewMockHTTPClient()
		geocoder := NewGeocoder(mockClient)
		if geocoder == nil {
			t.Errorf("Expected geocoder, but got nil")
		}
		if geocoder.client != mockClient {
			t.Errorf("Expected custom client to be set")
		}
	})
}
