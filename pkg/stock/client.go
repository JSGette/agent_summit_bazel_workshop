package stock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/JSGette/agent_summit_bazel_workshop/pkg/models"
)

// HTTPClient interface for dependency injection and testing
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// DefaultHTTPClient wraps the standard http.Client with proper headers
type DefaultHTTPClient struct{}

func (c *DefaultHTTPClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json,text/plain,*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	client := &http.Client{}
	return client.Do(req)
}

// Client handles stock API requests
type Client struct {
	httpClient HTTPClient
	baseURL    string
}

// NewClient creates a new stock client
func NewClient(httpClient HTTPClient) *Client {
	if httpClient == nil {
		httpClient = &DefaultHTTPClient{}
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    "https://query1.finance.yahoo.com/v7/finance/quote",
	}
}

// GetStockPrice fetches stock data for a given symbol
func (c *Client) GetStockPrice(symbol string) (*models.StockResponse, error) {
	if strings.TrimSpace(symbol) == "" {
		return nil, models.NewAPIError("Stock", "Symbol cannot be empty", 400)
	}

	// Normalize symbol to uppercase
	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	// Prepare URL with query parameters
	params := url.Values{}
	params.Add("symbols", symbol)

	requestURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	// Make the HTTP request
	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return nil, models.NewAPIError("Yahoo Finance", fmt.Sprintf("Failed to make request: %v", err), 500)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, models.NewAPIError("Yahoo Finance", fmt.Sprintf("API returned status %d", resp.StatusCode), resp.StatusCode)
	}

	// Parse the response
	var yahooResp models.YahooFinanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&yahooResp); err != nil {
		return nil, models.NewAPIError("Yahoo Finance", fmt.Sprintf("Failed to parse response: %v", err), 500)
	}

	// Convert to our standard format
	stockResp, err := models.ConvertYahooFinanceResponse(&yahooResp)
	if err != nil {
		return nil, err
	}

	return stockResp, nil
}

// GetDatadogStock is a convenience method to get Datadog (DDOG) stock price
func (c *Client) GetDatadogStock() (*models.StockResponse, error) {
	return c.GetStockPrice("DDOG")
}

// ValidateSymbol checks if a stock symbol is valid format
func (c *Client) ValidateSymbol(symbol string) error {
	symbol = strings.TrimSpace(symbol)

	if symbol == "" {
		return models.NewAPIError("Stock", "Symbol cannot be empty", 400)
	}

	if len(symbol) < 1 || len(symbol) > 5 {
		return models.NewAPIError("Stock", "Symbol must be 1-5 characters long", 400)
	}

	// Check if symbol contains only letters
	for _, char := range symbol {
		if !((char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')) {
			return models.NewAPIError("Stock", "Symbol must contain only letters", 400)
		}
	}

	return nil
}

// GetStockPriceWithValidation fetches stock data with input validation
func (c *Client) GetStockPriceWithValidation(symbol string) (*models.StockResponse, error) {
	if err := c.ValidateSymbol(symbol); err != nil {
		return nil, err
	}

	return c.GetStockPrice(symbol)
}
