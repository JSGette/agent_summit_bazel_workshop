package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/JSGette/agent_summit_bazel_workshop/pkg/models"
	"github.com/JSGette/agent_summit_bazel_workshop/pkg/stock"
	"github.com/JSGette/agent_summit_bazel_workshop/pkg/weather"
)

// Handler contains the services for handling HTTP requests
type Handler struct {
	weatherService *weather.Service
	stockService   *stock.Service
}

// NewHandler creates a new handler with the required services
func NewHandler(weatherService *weather.Service, stockService *stock.Service) *Handler {
	return &Handler{
		weatherService: weatherService,
		stockService:   stockService,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string    `json:"error"`
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Time    time.Time `json:"timestamp"`
}

// SuccessResponse represents a successful response wrapper
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Time    time.Time   `json:"timestamp"`
}

// writeErrorResponse writes an error response to the HTTP response writer
func (h *Handler) writeErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResp := ErrorResponse{
		Error:   err.Error(),
		Code:    statusCode,
		Message: "Request failed",
		Time:    time.Now(),
	}

	json.NewEncoder(w).Encode(errorResp)
	log.Printf("Error response: %v", err)
}

// writeSuccessResponse writes a successful response to the HTTP response writer
func (h *Handler) writeSuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	successResp := SuccessResponse{
		Success: true,
		Data:    data,
		Time:    time.Now(),
	}

	json.NewEncoder(w).Encode(successResp)
}

// GetWeather handles GET /weather?city=<city_name> requests
func (h *Handler) GetWeather(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, fmt.Errorf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	// Get city parameter from query string
	city := r.URL.Query().Get("city")
	if city == "" {
		h.writeErrorResponse(w, fmt.Errorf("missing required parameter 'city'"), http.StatusBadRequest)
		return
	}

	log.Printf("Weather request for city: %s", city)

	// Get weather data
	weatherData, err := h.weatherService.GetWeatherWithValidation(city)
	if err != nil {
		// Check if it's an API error to determine status code
		if apiErr, ok := err.(*models.APIError); ok {
			h.writeErrorResponse(w, err, apiErr.Code)
		} else {
			h.writeErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	h.writeSuccessResponse(w, weatherData)
	log.Printf("Weather request completed successfully for city: %s", city)
}

// GetDatadogStock handles GET /stock/datadog requests
func (h *Handler) GetDatadogStock(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, fmt.Errorf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Datadog stock price request")

	// Get Datadog stock data
	stockData, err := h.stockService.GetDatadogPrice()
	if err != nil {
		// Check if it's an API error to determine status code
		if apiErr, ok := err.(*models.APIError); ok {
			h.writeErrorResponse(w, err, apiErr.Code)
		} else {
			h.writeErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	h.writeSuccessResponse(w, stockData)
	log.Printf("Datadog stock request completed successfully")
}

// GetStock handles GET /stock?symbol=<symbol> requests (generic stock endpoint)
func (h *Handler) GetStock(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, fmt.Errorf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	// Get symbol parameter from query string
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		h.writeErrorResponse(w, fmt.Errorf("missing required parameter 'symbol'"), http.StatusBadRequest)
		return
	}

	log.Printf("Stock request for symbol: %s", symbol)

	// Get stock data
	stockData, err := h.stockService.GetCurrentPrice(symbol)
	if err != nil {
		// Check if it's an API error to determine status code
		if apiErr, ok := err.(*models.APIError); ok {
			h.writeErrorResponse(w, err, apiErr.Code)
		} else {
			h.writeErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	h.writeSuccessResponse(w, stockData)
	log.Printf("Stock request completed successfully for symbol: %s", symbol)
}

// HealthCheck handles GET /health requests
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, fmt.Errorf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	healthData := map[string]interface{}{
		"status":    "healthy",
		"service":   "weather-stock-api",
		"version":   "1.0.0",
		"timestamp": time.Now(),
		"uptime":    time.Since(startTime),
	}

	h.writeSuccessResponse(w, healthData)
}

// GetWeatherSummary handles GET /weather/summary?city=<city_name> requests
func (h *Handler) GetWeatherSummary(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, fmt.Errorf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	// Get city parameter from query string
	city := r.URL.Query().Get("city")
	if city == "" {
		h.writeErrorResponse(w, fmt.Errorf("missing required parameter 'city'"), http.StatusBadRequest)
		return
	}

	log.Printf("Weather summary request for city: %s", city)

	// Get weather summary
	summary, err := h.weatherService.GetWeatherSummary(city)
	if err != nil {
		// Check if it's an API error to determine status code
		if apiErr, ok := err.(*models.APIError); ok {
			h.writeErrorResponse(w, err, apiErr.Code)
		} else {
			h.writeErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	summaryData := map[string]interface{}{
		"city":    city,
		"summary": summary,
	}

	h.writeSuccessResponse(w, summaryData)
	log.Printf("Weather summary request completed successfully for city: %s", city)
}

// GetStockSummary handles GET /stock/summary?symbol=<symbol> requests
func (h *Handler) GetStockSummary(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, fmt.Errorf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	// Get symbol parameter from query string
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		h.writeErrorResponse(w, fmt.Errorf("missing required parameter 'symbol'"), http.StatusBadRequest)
		return
	}

	log.Printf("Stock summary request for symbol: %s", symbol)

	// Get stock summary
	summary, err := h.stockService.GetStockSummary(symbol)
	if err != nil {
		// Check if it's an API error to determine status code
		if apiErr, ok := err.(*models.APIError); ok {
			h.writeErrorResponse(w, err, apiErr.Code)
		} else {
			h.writeErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	summaryData := map[string]interface{}{
		"symbol":  symbol,
		"summary": summary,
	}

	h.writeSuccessResponse(w, summaryData)
	log.Printf("Stock summary request completed successfully for symbol: %s", symbol)
}

// Global variable to track server start time for uptime calculation
var startTime = time.Now()
