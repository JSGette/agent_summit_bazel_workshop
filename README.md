# Weather & Stock API

A simple Go web service that provides REST APIs for fetching weather information and stock prices, with **no API keys required**.

## Features

- 🌤️ **Weather API**: Get current weather for any city using Open-Meteo (free, no auth)
- 📈 **Stock API**: Get real-time stock prices using Yahoo Finance (free, no auth)
- 🎯 **Datadog Focus**: Special endpoint for Datadog (DDOG) stock price
- 🛡️ **Resilient Fallback**: Automatic demo mode when APIs are rate-limited or unavailable
- ✅ **Comprehensive Testing**: Full unit test coverage with mocked HTTP clients
- 🔧 **Production Ready**: Logging, middleware, graceful shutdown, error handling
- 🐳 **Zero Configuration**: No API keys, tokens, or external dependencies required

## Quick Start

### Build and Run
```bash
# Build the application
go build -o bin/weather-stock-api ./cmd

# Run the server
./bin/weather-stock-api

# Or run directly with Go
go run ./cmd
```

The server will start on `http://localhost:3000` by default.

### Docker (Optional)
```bash
# Build Docker image
docker build -t weather-stock-api .

# Run container
docker run -p 3000:3000 weather-stock-api
```

## API Endpoints

### Weather Endpoints

#### Get Weather Data
```bash
GET /weather?city=<city_name>

# Examples:
curl "http://localhost:3000/weather?city=Stuttgart"
curl "http://localhost:3000/weather?city=Berlin"
curl "http://localhost:3000/weather?city=New York"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "city": "Stuttgart",
    "country": "Germany",
    "temperature": 22.5,
    "condition": "partly_cloudy",
    "description": "Partly cloudy",
    "is_day": true,
    "coordinates": {
      "latitude": 48.7758,
      "longitude": 9.1829
    },
    "metadata": {
      "timestamp": "2024-01-15T14:30:00Z",
      "source": "Open-Meteo"
    }
  },
  "timestamp": "2024-01-15T14:30:15Z"
}
```

#### Get Weather Summary
```bash
GET /weather/summary?city=<city_name>

curl "http://localhost:3000/weather/summary?city=Stuttgart"
```

### Stock Endpoints

#### Get Datadog Stock Price
```bash
GET /stock/datadog

curl "http://localhost:3000/stock/datadog"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "symbol": "DDOG",
    "company_name": "Datadog, Inc.",
    "price": 125.67,
    "change": 2.34,
    "change_percent": 1.89,
    "previous_close": 123.33,
    "volume": 1234567,
    "market_cap": 40000000000,
    "market_state": "REGULAR",
    "currency": "USD",
    "metadata": {
      "timestamp": "2024-01-15T14:30:00Z",
      "source": "Yahoo Finance"
    }
  },
  "timestamp": "2024-01-15T14:30:15Z"
}
```

#### Get Any Stock Price
```bash
GET /stock?symbol=<symbol>

# Examples:
curl "http://localhost:3000/stock?symbol=DDOG"
curl "http://localhost:3000/stock?symbol=AAPL"
curl "http://localhost:3000/stock?symbol=GOOGL"
```

#### Get Stock Summary
```bash
GET /stock/summary?symbol=<symbol>

curl "http://localhost:3000/stock/summary?symbol=DDOG"
```

### Health Check
```bash
GET /health

curl "http://localhost:3000/health"
```

### API Information
```bash
GET /

curl "http://localhost:3000/"
```

## Configuration

The application supports configuration via environment variables and command-line flags:

### Environment Variables
- `HOST` - Server host (default: localhost)
- `PORT` - Server port (default: 3000)
- `READ_TIMEOUT` - HTTP read timeout (default: 10s)
- `WRITE_TIMEOUT` - HTTP write timeout (default: 10s)
- `IDLE_TIMEOUT` - HTTP idle timeout (default: 60s)

### Command Line Flags
```bash
./bin/weather-stock-api -help

# Custom configuration
./bin/weather-stock-api -host=0.0.0.0 -port=8080 -read-timeout=30s
```

### Examples
```bash
# Run on all interfaces, port 8080
PORT=8080 HOST=0.0.0.0 ./bin/weather-stock-api

# Run with custom timeouts
./bin/weather-stock-api -read-timeout=30s -write-timeout=30s
```

## Development

### Project Structure
```
agent_summit_bazel_workshop/
├── cmd/
│   └── main.go                 # Application entry point
├── pkg/
│   ├── models/                 # Data structures and types
│   │   ├── common.go          # Shared types and errors
│   │   ├── weather.go         # Weather response models
│   │   └── stock.go           # Stock response models
│   ├── weather/               # Weather service package
│   │   ├── client.go          # Open-Meteo API client
│   │   ├── geocode.go         # City to coordinates conversion
│   │   ├── service.go         # Business logic layer
│   │   └── *_test.go          # Unit tests
│   ├── stock/                 # Stock service package
│   │   ├── client.go          # Yahoo Finance API client
│   │   ├── service.go         # Business logic layer
│   │   └── *_test.go          # Unit tests
│   └── server/                # HTTP server package
│       ├── handlers.go        # HTTP request handlers
│       ├── middleware.go      # Middleware (logging, CORS, etc.)
│       ├── router.go          # Route definitions
│       └── server.go          # Server configuration
├── internal/
│   └── testutils/             # Testing utilities
│       ├── mocks.go           # HTTP client mocks
│       └── fixtures.go        # Test data fixtures
├── bin/                       # Compiled binaries
└── go.mod                     # Go module definition
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests for specific package
go test ./pkg/weather -v
go test ./pkg/stock -v

# Run tests with coverage
go test ./... -cover
```

### Building
```bash
# Build for current platform
go build -o bin/weather-stock-api ./cmd

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/weather-stock-api-linux ./cmd

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o bin/weather-stock-api.exe ./cmd
```

## Architecture

### Design Principles
- **Clean Architecture**: Separation of concerns with clear package boundaries
- **Dependency Injection**: HTTP clients are injected for easy testing
- **Interface-Based Design**: HTTP clients use interfaces for mockability  
- **Comprehensive Testing**: All business logic covered by unit tests
- **Error Handling**: Structured error responses with appropriate HTTP status codes
- **Middleware**: Logging, CORS, security headers, and panic recovery

### External APIs Used

#### Open-Meteo Weather API
- **URL**: https://api.open-meteo.com/v1/current
- **Authentication**: None required
- **Rate Limits**: Generous free tier
- **Features**: Current weather, geocoding, high accuracy

#### Yahoo Finance API
- **URL**: https://query1.finance.yahoo.com/v7/finance/quote
- **Authentication**: None required
- **Rate Limits**: Can be restrictive (429 errors)
- **Features**: Real-time stock data, market status, company info
- **Fallback**: Demo mode with realistic simulated data when rate-limited

## Example Workflow

```bash
# Start the server
./bin/weather-stock-api

# Check health
curl http://localhost:3000/health

# Get weather for Stuttgart
curl "http://localhost:3000/weather?city=Stuttgart"

# Get Datadog stock price
curl http://localhost:3000/stock/datadog

# Get weather summary
curl "http://localhost:3000/weather/summary?city=Stuttgart"

# Get stock summary
curl "http://localhost:3000/stock/summary?symbol=DDOG"
```

## Demo Mode & Fallback

When the Yahoo Finance API returns rate limiting errors (HTTP 429) or server errors (5xx), the stock service automatically falls back to **Demo Mode** which provides realistic simulated stock data.

### Demo Mode Features
- **Realistic Data**: Simulated prices based on real market patterns
- **Time-Aware**: Market state changes based on current time (open/closed/pre-market/after-hours)
- **Price Variation**: Prices vary slightly over time to simulate market movement
- **Multiple Stocks**: Supports DDOG, AAPL, GOOGL, MSFT, TSLA
- **Clear Indication**: Response includes "Demo Mode (Simulated Data)" in the source field

### Example Demo Response
```json
{
  "success": true,
  "data": {
    "symbol": "DDOG",
    "company_name": "Datadog, Inc.",
    "price": 127.45,
    "change": 1.95,
    "change_percent": 1.55,
    "market_state": "REGULAR",
    "metadata": {
      "timestamp": "2024-01-15T14:30:00Z",
      "source": "Demo Mode (Simulated Data)"
    }
  }
}
```

## Error Handling

The API returns structured error responses:

```json
{
  "error": "City 'InvalidCity' not found",
  "code": 404,
  "message": "Request failed",
  "timestamp": "2024-01-15T14:30:15Z"
}
```

Common HTTP status codes:
- `200` - Success
- `400` - Bad Request (missing parameters, invalid input)
- `404` - Not Found (city not found, stock symbol not found)
- `500` - Internal Server Error (API failures, network issues)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass: `go test ./...`
5. Build successfully: `go build ./cmd`
6. Submit a pull request

## License

MIT License - see LICENSE file for details.

---

**Built with Go 1.24.4** • **No API keys required** • **Ready for production**
