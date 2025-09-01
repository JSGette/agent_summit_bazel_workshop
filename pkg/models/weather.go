package models

import "time"

// WeatherCondition represents different weather states
type WeatherCondition string

const (
	Clear        WeatherCondition = "clear"
	PartlyCloudy WeatherCondition = "partly_cloudy"
	Cloudy       WeatherCondition = "cloudy"
	Overcast     WeatherCondition = "overcast"
	Fog          WeatherCondition = "fog"
	Drizzle      WeatherCondition = "drizzle"
	Rain         WeatherCondition = "rain"
	Snow         WeatherCondition = "snow"
	Thunderstorm WeatherCondition = "thunderstorm"
	Unknown      WeatherCondition = "unknown"
)

// WeatherResponse represents the standardized weather response
type WeatherResponse struct {
	City        string           `json:"city"`
	Country     string           `json:"country"`
	Temperature float64          `json:"temperature"`
	Condition   WeatherCondition `json:"condition"`
	Description string           `json:"description"`
	IsDay       bool             `json:"is_day"`
	Coordinates Coordinates      `json:"coordinates"`
	Metadata    ResponseMetadata `json:"metadata"`
}

// OpenMeteoResponse represents the raw response from Open-Meteo API
type OpenMeteoResponse struct {
	Current struct {
		Time          string  `json:"time"`
		Temperature2m float64 `json:"temperature_2m"`
		WeatherCode   int     `json:"weather_code"`
		IsDay         int     `json:"is_day"`
	} `json:"current"`
	CurrentUnits struct {
		Temperature2m string `json:"temperature_2m"`
	} `json:"current_units"`
}

// WeatherCodeMap maps Open-Meteo weather codes to our conditions
var WeatherCodeMap = map[int]struct {
	Condition   WeatherCondition
	Description string
}{
	0:  {Clear, "Clear sky"},
	1:  {PartlyCloudy, "Mainly clear"},
	2:  {PartlyCloudy, "Partly cloudy"},
	3:  {Cloudy, "Overcast"},
	45: {Fog, "Fog"},
	48: {Fog, "Depositing rime fog"},
	51: {Drizzle, "Light drizzle"},
	53: {Drizzle, "Moderate drizzle"},
	55: {Drizzle, "Dense drizzle"},
	61: {Rain, "Slight rain"},
	63: {Rain, "Moderate rain"},
	65: {Rain, "Heavy rain"},
	71: {Snow, "Slight snow fall"},
	73: {Snow, "Moderate snow fall"},
	75: {Snow, "Heavy snow fall"},
	95: {Thunderstorm, "Thunderstorm"},
	96: {Thunderstorm, "Thunderstorm with slight hail"},
	99: {Thunderstorm, "Thunderstorm with heavy hail"},
}

// GetWeatherCondition converts Open-Meteo weather code to our condition
func GetWeatherCondition(code int) (WeatherCondition, string) {
	if weather, exists := WeatherCodeMap[code]; exists {
		return weather.Condition, weather.Description
	}
	return Unknown, "Unknown weather condition"
}

// ConvertOpenMeteoResponse converts Open-Meteo API response to our standard format
func ConvertOpenMeteoResponse(response *OpenMeteoResponse, city, country string, coords Coordinates) *WeatherResponse {
	condition, description := GetWeatherCondition(response.Current.WeatherCode)

	// Parse time
	timestamp, _ := time.Parse("2006-01-02T15:04", response.Current.Time)

	return &WeatherResponse{
		City:        city,
		Country:     country,
		Temperature: response.Current.Temperature2m,
		Condition:   condition,
		Description: description,
		IsDay:       response.Current.IsDay == 1,
		Coordinates: coords,
		Metadata: ResponseMetadata{
			Timestamp: timestamp,
			Source:    "Open-Meteo",
		},
	}
}
