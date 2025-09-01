package stock

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/JSGette/agent_summit_bazel_workshop/pkg/models"
)

// Service provides high-level stock operations with caching and logging
type Service struct {
	client      *Client
	lastRequest time.Time
	mutex       sync.Mutex
}

// NewService creates a new stock service
func NewService(httpClient HTTPClient) *Service {
	return &Service{
		client: NewClient(httpClient),
	}
}

// rateLimitDelay enforces a minimum delay between API requests
func (s *Service) rateLimitDelay() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	const minDelay = 2 * time.Second // 2 seconds between requests
	timeSinceLastRequest := time.Since(s.lastRequest)

	if timeSinceLastRequest < minDelay {
		sleepTime := minDelay - timeSinceLastRequest
		log.Printf("Rate limiting: sleeping for %v", sleepTime)
		time.Sleep(sleepTime)
	}

	s.lastRequest = time.Now()
}

// GetCurrentPrice fetches current stock price for a symbol with enhanced error handling
func (s *Service) GetCurrentPrice(symbol string) (*models.StockResponse, error) {
	start := time.Now()

	log.Printf("Fetching stock price for symbol: %s", symbol)

	// Apply rate limiting
	s.rateLimitDelay()

	stock, err := s.client.GetStockPriceWithValidation(symbol)
	if err != nil {
		log.Printf("Error fetching stock price for %s: %v", symbol, err)

		// Check if it's a rate limit error (429), auth error (401/403), or server error (5xx) - fall back to demo mode
		if apiErr, ok := err.(*models.APIError); ok && (apiErr.Code == 401 || apiErr.Code == 403 || apiErr.Code == 429 || apiErr.Code >= 500) {
			log.Printf("API error %d, falling back to demo mode for %s", apiErr.Code, symbol)
			demoStock, demoErr := GetDemoStock(symbol)
			if demoErr != nil {
				log.Printf("Demo mode also failed for %s: %v", symbol, demoErr)
				return nil, err // Return original error
			}
			log.Printf("Successfully returned demo data for %s", symbol)
			return demoStock, nil
		}

		return nil, err
	}

	duration := time.Since(start)
	log.Printf("Successfully fetched stock price for %s in %v", symbol, duration)

	return stock, nil
}

// GetDatadogPrice is a convenience method to get Datadog stock price
func (s *Service) GetDatadogPrice() (*models.StockResponse, error) {
	return s.GetCurrentPrice("DDOG")
}

// GetStockSummary returns a human-readable stock summary
func (s *Service) GetStockSummary(symbol string) (string, error) {
	stock, err := s.GetCurrentPrice(symbol)
	if err != nil {
		return "", err
	}

	direction := "unchanged"
	changeIcon := "→"

	if stock.Change > 0 {
		direction = "up"
		changeIcon = "↗"
	} else if stock.Change < 0 {
		direction = "down"
		changeIcon = "↘"
	}

	marketStateText := ""
	switch stock.MarketState {
	case models.MarketStateRegular:
		marketStateText = "Market Open"
	case models.MarketStatePremarket:
		marketStateText = "Pre-Market"
	case models.MarketStatePostmarket:
		marketStateText = "After Hours"
	case models.MarketStateClosed:
		marketStateText = "Market Closed"
	}

	summary := fmt.Sprintf(
		"%s (%s): $%.2f %s %.2f (%.2f%%) - %s. %s. Last updated: %s",
		stock.CompanyName,
		stock.Symbol,
		stock.Price,
		changeIcon,
		stock.Change,
		stock.ChangePercent,
		direction,
		marketStateText,
		stock.Metadata.Timestamp.Format("15:04 MST"),
	)

	return summary, nil
}

// GetDatadogSummary returns a formatted summary for Datadog stock
func (s *Service) GetDatadogSummary() (string, error) {
	return s.GetStockSummary("DDOG")
}

// IsMarketOpen checks if the market is currently open based on the stock data
func (s *Service) IsMarketOpen(symbol string) (bool, error) {
	stock, err := s.GetCurrentPrice(symbol)
	if err != nil {
		return false, err
	}

	return stock.MarketState == models.MarketStateRegular, nil
}

// GetPriceChange returns formatted price change information
func (s *Service) GetPriceChange(symbol string) (string, error) {
	stock, err := s.GetCurrentPrice(symbol)
	if err != nil {
		return "", err
	}

	sign := ""
	if stock.Change > 0 {
		sign = "+"
	}

	return fmt.Sprintf("%s%.2f (%.2f%%)", sign, stock.Change, stock.ChangePercent), nil
}

// ValidateAndNormalizeSymbol validates and normalizes a stock symbol
func (s *Service) ValidateAndNormalizeSymbol(symbol string) (string, error) {
	if err := s.client.ValidateSymbol(symbol); err != nil {
		return "", err
	}

	// Return normalized (uppercase, trimmed) symbol
	return fmt.Sprintf("%s", symbol), nil
}
