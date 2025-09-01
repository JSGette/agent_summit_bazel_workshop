package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/JSGette/agent_summit_bazel_workshop/pkg/models"
)

// GeocodeResponse represents the response from Open-Meteo geocoding API
type GeocodeResponse struct {
	Results []struct {
		Name        string  `json:"name"`
		Country     string  `json:"country"`
		CountryCode string  `json:"country_code"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		Admin1      string  `json:"admin1,omitempty"`
	} `json:"results"`
}

// HTTPClient interface for dependency injection and testing
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// DefaultHTTPClient wraps the standard http.Client
type DefaultHTTPClient struct{}

func (c *DefaultHTTPClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

// Geocoder handles city name to coordinates conversion
type Geocoder struct {
	client  HTTPClient
	baseURL string
}

// NewGeocoder creates a new geocoder instance
func NewGeocoder(client HTTPClient) *Geocoder {
	if client == nil {
		client = &DefaultHTTPClient{}
	}
	return &Geocoder{
		client:  client,
		baseURL: "https://geocoding-api.open-meteo.com/v1/search",
	}
}

// GetCoordinates converts a city name to coordinates using Open-Meteo geocoding API
func (g *Geocoder) GetCoordinates(city string) (*models.Coordinates, string, error) {
	if strings.TrimSpace(city) == "" {
		return nil, "", models.NewAPIError("Geocoding", "City name cannot be empty", 400)
	}

	// Prepare the URL with query parameters
	params := url.Values{}
	params.Add("name", city)
	params.Add("count", "1")
	params.Add("language", "en")
	params.Add("format", "json")

	requestURL := fmt.Sprintf("%s?%s", g.baseURL, params.Encode())

	// Make the HTTP request
	resp, err := g.client.Get(requestURL)
	if err != nil {
		return nil, "", models.NewAPIError("Geocoding", fmt.Sprintf("Failed to make request: %v", err), 500)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", models.NewAPIError("Geocoding", fmt.Sprintf("API returned status %d", resp.StatusCode), resp.StatusCode)
	}

	// Parse the response
	var geocodeResp GeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&geocodeResp); err != nil {
		return nil, "", models.NewAPIError("Geocoding", fmt.Sprintf("Failed to parse response: %v", err), 500)
	}

	// Check if we got any results
	if len(geocodeResp.Results) == 0 {
		return nil, "", models.NewAPIError("Geocoding", fmt.Sprintf("City '%s' not found", city), 404)
	}

	result := geocodeResp.Results[0]
	coords := &models.Coordinates{
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
	}

	return coords, result.Country, nil
}

// CityCoordinates is a simple in-memory cache for common cities
var CityCoordinates = map[string]struct {
	Coords  models.Coordinates
	Country string
}{
	"stuttgart": {
		Coords:  models.Coordinates{Latitude: 48.7758, Longitude: 9.1829},
		Country: "Germany",
	},
	"berlin": {
		Coords:  models.Coordinates{Latitude: 52.5200, Longitude: 13.4050},
		Country: "Germany",
	},
	"munich": {
		Coords:  models.Coordinates{Latitude: 48.1351, Longitude: 11.5820},
		Country: "Germany",
	},
	"london": {
		Coords:  models.Coordinates{Latitude: 51.5074, Longitude: -0.1278},
		Country: "United Kingdom",
	},
	"paris": {
		Coords:  models.Coordinates{Latitude: 48.8566, Longitude: 2.3522},
		Country: "France",
	},
	"new york": {
		Coords:  models.Coordinates{Latitude: 40.7128, Longitude: -74.0060},
		Country: "United States",
	},
}

// GetCoordinatesWithCache tries cache first, then falls back to API
func (g *Geocoder) GetCoordinatesWithCache(city string) (*models.Coordinates, string, error) {
	cityLower := strings.ToLower(strings.TrimSpace(city))

	// Check cache first
	if cached, exists := CityCoordinates[cityLower]; exists {
		return &cached.Coords, cached.Country, nil
	}

	// Fall back to API
	return g.GetCoordinates(city)
}
