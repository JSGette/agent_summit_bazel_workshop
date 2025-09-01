package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/JSGette/agent_summit_bazel_workshop/pkg/server"
	"github.com/JSGette/agent_summit_bazel_workshop/pkg/stock"
	"github.com/JSGette/agent_summit_bazel_workshop/pkg/weather"
)

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Weather & Stock API service...")

	// Parse command line flags
	var (
		host         = flag.String("host", getEnv("HOST", "localhost"), "Server host")
		port         = flag.Int("port", getEnvInt("PORT", 3000), "Server port")
		readTimeout  = flag.Duration("read-timeout", getEnvDuration("READ_TIMEOUT", "10s"), "HTTP read timeout")
		writeTimeout = flag.Duration("write-timeout", getEnvDuration("WRITE_TIMEOUT", "10s"), "HTTP write timeout")
		idleTimeout  = flag.Duration("idle-timeout", getEnvDuration("IDLE_TIMEOUT", "60s"), "HTTP idle timeout")
		showHelp     = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *showHelp {
		showUsage()
		return
	}

	// Create server configuration
	config := &server.Config{
		Host:         *host,
		Port:         *port,
		ReadTimeout:  *readTimeout,
		WriteTimeout: *writeTimeout,
		IdleTimeout:  *idleTimeout,
	}

	// Initialize services
	log.Println("Initializing services...")

	// Create HTTP client (nil will use default)
	var httpClient interface {
		Get(url string) (*http.Response, error)
	}
	// httpClient = nil // Use default HTTP client

	// Initialize weather service
	weatherService := weather.NewService(httpClient)
	log.Println("Weather service initialized")

	// Initialize stock service
	stockService := stock.NewService(httpClient)
	log.Println("Stock service initialized")

	// Create and configure server
	srv := server.NewServer(config, weatherService, stockService)
	log.Printf("Server created and configured to run on %s:%d", config.Host, config.Port)

	// Start server with graceful shutdown
	log.Println("Starting server...")
	if err := srv.StartWithGracefulShutdown(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	log.Println("Server shutdown complete")
}

// showUsage displays usage information
func showUsage() {
	log.Println("Weather & Stock API Server")
	log.Println("")
	log.Println("This service provides REST API endpoints for:")
	log.Println("  - Current weather information for cities")
	log.Println("  - Real-time stock prices (including Datadog)")
	log.Println("")
	log.Println("Environment Variables:")
	log.Println("  HOST         - Server host (default: localhost)")
	log.Println("  PORT         - Server port (default: 3000)")
	log.Println("  READ_TIMEOUT - HTTP read timeout (default: 10s)")
	log.Println("  WRITE_TIMEOUT- HTTP write timeout (default: 10s)")
	log.Println("  IDLE_TIMEOUT - HTTP idle timeout (default: 60s)")
	log.Println("")
	log.Println("Command Line Flags:")
	flag.PrintDefaults()
	log.Println("")
	log.Println("API Endpoints:")
	log.Println("  GET /health                     - Health check")
	log.Println("  GET /weather?city=<name>        - Get weather for city")
	log.Println("  GET /weather/summary?city=<name>- Get weather summary")
	log.Println("  GET /stock?symbol=<symbol>      - Get stock price")
	log.Println("  GET /stock/datadog              - Get Datadog stock price")
	log.Println("  GET /stock/summary?symbol=<sym> - Get stock summary")
	log.Println("")
	log.Println("Examples:")
	log.Println("  curl http://localhost:3000/weather?city=Stuttgart")
	log.Println("  curl http://localhost:3000/stock/datadog")
	log.Println("  curl http://localhost:3000/health")
}

// getEnv returns environment variable value or default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt returns environment variable as int or default
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		log.Printf("Warning: Invalid integer value for %s: %s, using default %d", key, value, defaultValue)
	}
	return defaultValue
}

// getEnvDuration returns environment variable as duration or default
func getEnvDuration(key, defaultValue string) time.Duration {
	value := getEnv(key, defaultValue)
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	log.Printf("Warning: Invalid duration value for %s: %s, using default %s", key, value, defaultValue)
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return 10 * time.Second // fallback
}

// Version information (could be set during build)
var (
	Version   = "1.0.0"
	BuildTime = "development"
	GitCommit = "unknown"
)

func init() {
	log.Printf("Weather & Stock API v%s (built: %s, commit: %s)", Version, BuildTime, GitCommit)
}
