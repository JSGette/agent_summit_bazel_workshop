package stock

import (
	"errors"
	"strings"
	"testing"

	"github.com/JSGette/agent_summit_bazel_workshop/internal/testutils"
)

func TestClient_GetStockPrice(t *testing.T) {
	tests := []struct {
		name           string
		symbol         string
		mockResponse   string
		mockStatusCode int
		mockError      error
		wantError      bool
		wantSymbol     string
		wantPrice      float64
		wantChange     float64
	}{
		{
			name:           "successful stock request",
			symbol:         "DDOG",
			mockResponse:   testutils.YahooFinanceStockResponse,
			mockStatusCode: 200,
			wantError:      false,
			wantSymbol:     "DDOG",
			wantPrice:      125.67,
			wantChange:     2.34,
		},
		{
			name:           "symbol case insensitive",
			symbol:         "ddog",
			mockResponse:   testutils.YahooFinanceStockResponse,
			mockStatusCode: 200,
			wantError:      false,
			wantSymbol:     "DDOG",
			wantPrice:      125.67,
			wantChange:     2.34,
		},
		{
			name:           "symbol with whitespace",
			symbol:         "  DDOG  ",
			mockResponse:   testutils.YahooFinanceStockResponse,
			mockStatusCode: 200,
			wantError:      false,
			wantSymbol:     "DDOG",
			wantPrice:      125.67,
			wantChange:     2.34,
		},
		{
			name:      "empty symbol",
			symbol:    "",
			wantError: true,
		},
		{
			name:      "whitespace only symbol",
			symbol:    "   ",
			wantError: true,
		},
		{
			name:           "stock not found",
			symbol:         "INVALID",
			mockResponse:   testutils.YahooFinanceStockNotFound,
			mockStatusCode: 200,
			wantError:      true,
		},
		{
			name:           "API returns 500 error",
			symbol:         "DDOG",
			mockResponse:   testutils.APIErrorResponse,
			mockStatusCode: 500,
			wantError:      true,
		},
		{
			name:      "network error",
			symbol:    "DDOG",
			mockError: errors.New("network error"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := testutils.NewMockHTTPClient()
			client := NewClient(mockClient)

			if tt.symbol != "" && strings.TrimSpace(tt.symbol) != "" {
				// Prepare expected URL - symbol should be normalized to uppercase
				expectedSymbol := strings.ToUpper(strings.TrimSpace(tt.symbol))
				expectedURL := "https://query1.finance.yahoo.com/v7/finance/quote?symbols=" + expectedSymbol

				if tt.mockError != nil {
					mockClient.AddError(expectedURL, tt.mockError)
				} else {
					mockClient.AddResponse(expectedURL, tt.mockStatusCode, tt.mockResponse)
				}
			}

			result, err := client.GetStockPrice(tt.symbol)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("Expected result, but got nil")
				return
			}

			if result.Symbol != tt.wantSymbol {
				t.Errorf("Expected symbol %v, got %v", tt.wantSymbol, result.Symbol)
			}

			if result.Price != tt.wantPrice {
				t.Errorf("Expected price %v, got %v", tt.wantPrice, result.Price)
			}

			if result.Change != tt.wantChange {
				t.Errorf("Expected change %v, got %v", tt.wantChange, result.Change)
			}
		})
	}
}

func TestClient_GetDatadogStock(t *testing.T) {
	mockClient := testutils.NewMockHTTPClient()
	client := NewClient(mockClient)

	expectedURL := "https://query1.finance.yahoo.com/v7/finance/quote?symbols=DDOG"
	mockClient.AddResponse(expectedURL, 200, testutils.YahooFinanceStockResponse)

	result, err := client.GetDatadogStock()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if result == nil {
		t.Errorf("Expected result, but got nil")
		return
	}

	if result.Symbol != "DDOG" {
		t.Errorf("Expected symbol DDOG, got %v", result.Symbol)
	}
}

func TestClient_ValidateSymbol(t *testing.T) {
	tests := []struct {
		name      string
		symbol    string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid symbol",
			symbol:    "DDOG",
			wantError: false,
		},
		{
			name:      "valid single character",
			symbol:    "A",
			wantError: false,
		},
		{
			name:      "valid 5 characters",
			symbol:    "AAPLE",
			wantError: false,
		},
		{
			name:      "valid lowercase",
			symbol:    "ddog",
			wantError: false,
		},
		{
			name:      "empty symbol",
			symbol:    "",
			wantError: true,
			errorMsg:  "Symbol cannot be empty",
		},
		{
			name:      "whitespace only",
			symbol:    "   ",
			wantError: true,
			errorMsg:  "Symbol cannot be empty",
		},
		{
			name:      "too long symbol",
			symbol:    "TOOLONG",
			wantError: true,
			errorMsg:  "1-5 characters long",
		},
		{
			name:      "symbol with numbers",
			symbol:    "DD0G",
			wantError: true,
			errorMsg:  "contain only letters",
		},
		{
			name:      "symbol with special characters",
			symbol:    "DD-G",
			wantError: true,
			errorMsg:  "contain only letters",
		},
	}

	client := NewClient(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ValidateSymbol(tt.symbol)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error, but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', got: %v", tt.errorMsg, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestClient_GetStockPriceWithValidation(t *testing.T) {
	tests := []struct {
		name      string
		symbol    string
		wantError bool
	}{
		{
			name:      "valid symbol",
			symbol:    "DDOG",
			wantError: false,
		},
		{
			name:      "invalid symbol - empty",
			symbol:    "",
			wantError: true,
		},
		{
			name:      "invalid symbol - too long",
			symbol:    "TOOLONG",
			wantError: true,
		},
		{
			name:      "invalid symbol - with numbers",
			symbol:    "DD0G",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := testutils.NewMockHTTPClient()
			client := NewClient(mockClient)

			// Mock successful API response for valid symbols
			if !tt.wantError {
				expectedURL := "https://query1.finance.yahoo.com/v7/finance/quote?symbols=" + strings.ToUpper(tt.symbol)
				mockClient.AddResponse(expectedURL, 200, testutils.YahooFinanceStockResponse)
			}

			_, err := client.GetStockPriceWithValidation(tt.symbol)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	t.Run("with nil client", func(t *testing.T) {
		client := NewClient(nil)
		if client == nil {
			t.Errorf("Expected client, but got nil")
		}
		if client.httpClient == nil {
			t.Errorf("Expected default HTTP client to be set")
		}
	})

	t.Run("with custom client", func(t *testing.T) {
		mockClient := testutils.NewMockHTTPClient()
		client := NewClient(mockClient)
		if client == nil {
			t.Errorf("Expected client, but got nil")
		}
		if client.httpClient != mockClient {
			t.Errorf("Expected custom client to be set")
		}
	})
}
