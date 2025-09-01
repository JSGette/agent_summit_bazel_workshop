package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JSGette/agent_summit_bazel_workshop/pkg/stock"
	"github.com/JSGette/agent_summit_bazel_workshop/pkg/weather"
)

// Server represents the HTTP server
type Server struct {
	httpServer     *http.Server
	weatherService *weather.Service
	stockService   *stock.Service
	router         *Router
}

// Config holds server configuration
type Config struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DefaultConfig returns default server configuration
func DefaultConfig() *Config {
	return &Config{
		Host:         "localhost",
		Port:         3000,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// NewServer creates a new server instance
func NewServer(config *Config, weatherService *weather.Service, stockService *stock.Service) *Server {
	if config == nil {
		config = DefaultConfig()
	}

	router := NewRouter(weatherService, stockService)

	server := &Server{
		weatherService: weatherService,
		stockService:   stockService,
		router:         router,
	}

	server.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler:      router.GetHandler(),
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	return server
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("Starting server on %s", s.httpServer.Addr)
	log.Printf("Server configuration:")
	log.Printf("  Read timeout: %v", s.httpServer.ReadTimeout)
	log.Printf("  Write timeout: %v", s.httpServer.WriteTimeout)
	log.Printf("  Idle timeout: %v", s.httpServer.IdleTimeout)

	// Print available endpoints
	s.printAvailableEndpoints()

	return s.httpServer.ListenAndServe()
}

// StartWithGracefulShutdown starts the server with graceful shutdown support
func (s *Server) StartWithGracefulShutdown() error {
	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := s.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for signal
	sig := <-sigChan
	log.Printf("Received signal: %v", sig)

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	return s.Shutdown(ctx)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// printAvailableEndpoints prints all available API endpoints
func (s *Server) printAvailableEndpoints() {
	baseURL := fmt.Sprintf("http://%s", s.httpServer.Addr)

	log.Println("Available endpoints:")
	log.Printf("  GET %s/                    - API information", baseURL)
	log.Printf("  GET %s/health              - Health check", baseURL)
	log.Printf("  GET %s/weather?city=<name> - Get weather (example: ?city=Stuttgart)", baseURL)
	log.Printf("  GET %s/weather/summary?city=<name> - Get weather summary", baseURL)
	log.Printf("  GET %s/stock?symbol=<sym>  - Get stock price (example: ?symbol=DDOG)", baseURL)
	log.Printf("  GET %s/stock/datadog       - Get Datadog stock price", baseURL)
	log.Printf("  GET %s/stock/summary?symbol=<sym> - Get stock summary", baseURL)
	log.Println()
}

// GetAddr returns the server address
func (s *Server) GetAddr() string {
	return s.httpServer.Addr
}

// IsRunning checks if the server is running (simplified check)
func (s *Server) IsRunning() bool {
	return s.httpServer != nil
}
