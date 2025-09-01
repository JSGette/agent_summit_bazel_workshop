package testutils

// Weather API Response Fixtures

// OpenMeteoWeatherResponse is a sample response from Open-Meteo API
const OpenMeteoWeatherResponse = `{
  "current": {
    "time": "2024-01-15T14:00",
    "temperature_2m": 22.5,
    "weather_code": 3,
    "is_day": 1
  },
  "current_units": {
    "temperature_2m": "°C"
  }
}`

// OpenMeteoGeocodeResponse is a sample response from Open-Meteo Geocoding API
const OpenMeteoGeocodeResponse = `{
  "results": [
    {
      "name": "Stuttgart",
      "country": "Germany",
      "country_code": "DE",
      "latitude": 48.7758,
      "longitude": 9.1829,
      "admin1": "Baden-Württemberg"
    }
  ]
}`

// OpenMeteoGeocodeNotFound is a response when city is not found
const OpenMeteoGeocodeNotFound = `{
  "results": []
}`

// Stock API Response Fixtures

// YahooFinanceStockResponse is a sample response from Yahoo Finance API
const YahooFinanceStockResponse = `{
  "quoteResponse": {
    "result": [
      {
        "symbol": "DDOG",
        "shortName": "Datadog Inc",
        "longName": "Datadog, Inc.",
        "regularMarketPrice": 125.67,
        "regularMarketChange": 2.34,
        "regularMarketChangePercent": 1.89,
        "regularMarketPreviousClose": 123.33,
        "regularMarketVolume": 1234567,
        "marketCap": 40000000000,
        "currency": "USD",
        "marketState": "REGULAR",
        "regularMarketTime": 1705327200
      }
    ],
    "error": null
  }
}`

// YahooFinanceStockNotFound is a response when stock symbol is not found
const YahooFinanceStockNotFound = `{
  "quoteResponse": {
    "result": [],
    "error": null
  }
}`

// YahooFinanceMarketClosed is a response when market is closed
const YahooFinanceMarketClosed = `{
  "quoteResponse": {
    "result": [
      {
        "symbol": "DDOG",
        "shortName": "Datadog Inc",
        "longName": "Datadog, Inc.",
        "regularMarketPrice": 125.67,
        "regularMarketChange": -1.23,
        "regularMarketChangePercent": -0.97,
        "regularMarketPreviousClose": 126.90,
        "regularMarketVolume": 987654,
        "marketCap": 40000000000,
        "currency": "USD",
        "marketState": "CLOSED",
        "regularMarketTime": 1705327200
      }
    ],
    "error": null
  }
}`

// Error Response Fixtures

// APIErrorResponse is a generic API error response
const APIErrorResponse = `{
  "error": {
    "code": 500,
    "message": "Internal server error"
  }
}`

// RateLimitErrorResponse simulates a rate limit error
const RateLimitErrorResponse = `{
  "error": {
    "code": 429,
    "message": "Too many requests"
  }
}`
