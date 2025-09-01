package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/JSGette/agent_summit_bazel_workshop/pkg/models"
)

// Client handles weather API requests
type Client struct {
	httpClient HTTPClient
	geocoder   *Geocoder
	baseURL    string
}

// NewClient creates a new weather client
func NewClient(httpClient HTTPClient) *Client {
	if httpClient == nil {
		httpClient = &DefaultHTTPClient{}
	}

	return &Client{
		httpClient: httpClient,
		geocoder:   NewGeocoder(httpClient),
		baseURL:    "https://api.open-meteo.com/v1/forecast",
	}
}

// GetWeatherByCity fetches weather data for a given city name
func (c *Client) GetWeatherByCity(city string) (*models.WeatherResponse, error) {
	// Get coordinates for the city
	coords, country, err := c.geocoder.GetCoordinatesWithCache(city)
	if err != nil {
		return nil, err
	}

	// Get weather data using coordinates
	return c.GetWeatherByCoordinates(coords.Latitude, coords.Longitude, city, country)
}

// GetWeatherByCoordinates fetches weather data for given coordinates
func (c *Client) GetWeatherByCoordinates(lat, lon float64, city, country string) (*models.WeatherResponse, error) {
	// Prepare URL with query parameters
	params := url.Values{}
	params.Add("latitude", fmt.Sprintf("%.4f", lat))
	params.Add("longitude", fmt.Sprintf("%.4f", lon))
	params.Add("current", "temperature_2m,weather_code,is_day")
	params.Add("timezone", "auto")

	requestURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	// Make the HTTP request
	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return nil, models.NewAPIError("Open-Meteo", fmt.Sprintf("Failed to make request: %v", err), 500)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, models.NewAPIError("Open-Meteo", fmt.Sprintf("API returned status %d", resp.StatusCode), resp.StatusCode)
	}

	// Parse the response
	var openMeteoResp models.OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&openMeteoResp); err != nil {
		return nil, models.NewAPIError("Open-Meteo", fmt.Sprintf("Failed to parse response: %v", err), 500)
	}

	// Convert to our standard format
	coords := models.Coordinates{Latitude: lat, Longitude: lon}
	weatherResp := models.ConvertOpenMeteoResponse(&openMeteoResp, city, country, coords)

	return weatherResp, nil
}

// GetWeather is a convenience method that handles both city names and coordinates
func (c *Client) GetWeather(location string) (*models.WeatherResponse, error) {
	if location == "" {
		return nil, models.NewAPIError("Weather", "Location cannot be empty", 400)
	}

	// For now, treat all inputs as city names
	// In the future, we could add support for "lat,lon" format
	return c.GetWeatherByCity(location)
}
