package server

import (
	"net/http"

	"github.com/JSGette/agent_summit_bazel_workshop/pkg/stock"
	"github.com/JSGette/agent_summit_bazel_workshop/pkg/weather"
)

// Router handles HTTP routing
type Router struct {
	handler *Handler
	mux     *http.ServeMux
}

// NewRouter creates a new router with all routes configured
func NewRouter(weatherService *weather.Service, stockService *stock.Service) *Router {
	handler := NewHandler(weatherService, stockService)
	mux := http.NewServeMux()

	router := &Router{
		handler: handler,
		mux:     mux,
	}

	router.setupRoutes()
	return router
}

// setupRoutes configures all the HTTP routes
func (router *Router) setupRoutes() {
	// Health check endpoint
	router.mux.HandleFunc("/health", router.handler.HealthCheck)

	// Weather endpoints
	router.mux.HandleFunc("/weather", router.handler.GetWeather)
	router.mux.HandleFunc("/weather/summary", router.handler.GetWeatherSummary)

	// Stock endpoints
	router.mux.HandleFunc("/stock", router.handler.GetStock)
	router.mux.HandleFunc("/stock/datadog", router.handler.GetDatadogStock)
	router.mux.HandleFunc("/stock/summary", router.handler.GetStockSummary)

	// Add a root endpoint for basic info
	router.mux.HandleFunc("/", router.rootHandler)
}

// rootHandler provides basic API information
func (router *Router) rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		router.handler.writeErrorResponse(w, http.ErrNotSupported, http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	apiInfo := map[string]interface{}{
		"service":     "Weather & Stock API",
		"version":     "1.0.0",
		"description": "A simple API to get weather information and stock prices",
		"endpoints": map[string]interface{}{
			"health": map[string]string{
				"method":      "GET",
				"path":        "/health",
				"description": "Health check endpoint",
			},
			"weather": map[string]string{
				"method":      "GET",
				"path":        "/weather?city=<city_name>",
				"description": "Get current weather for a city",
				"example":     "/weather?city=Stuttgart",
			},
			"weather_summary": map[string]string{
				"method":      "GET",
				"path":        "/weather/summary?city=<city_name>",
				"description": "Get weather summary for a city",
				"example":     "/weather/summary?city=Stuttgart",
			},
			"stock": map[string]string{
				"method":      "GET",
				"path":        "/stock?symbol=<symbol>",
				"description": "Get current stock price for a symbol",
				"example":     "/stock?symbol=DDOG",
			},
			"datadog_stock": map[string]string{
				"method":      "GET",
				"path":        "/stock/datadog",
				"description": "Get current Datadog stock price",
			},
			"stock_summary": map[string]string{
				"method":      "GET",
				"path":        "/stock/summary?symbol=<symbol>",
				"description": "Get stock summary for a symbol",
				"example":     "/stock/summary?symbol=DDOG",
			},
		},
	}

	router.handler.writeSuccessResponse(w, apiInfo)
}

// ServeHTTP implements the http.Handler interface
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

// GetHandler returns the configured HTTP handler with middleware
func (router *Router) GetHandler() http.Handler {
	// Apply middleware in reverse order (last applied is executed first)
	var handler http.Handler = router.mux
	handler = SecurityMiddleware(handler)
	handler = ContentTypeMiddleware(handler)
	handler = CORSMiddleware(handler)
	handler = RecoveryMiddleware(handler)
	handler = LoggingMiddleware(handler)

	return handler
}
