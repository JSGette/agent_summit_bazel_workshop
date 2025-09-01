package stock

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/JSGette/agent_summit_bazel_workshop/internal/testutils"
	"github.com/JSGette/agent_summit_bazel_workshop/pkg/models"
)

func TestService_GetCurrentPrice(t *testing.T) {
	tests := []struct {
		name           string
		symbol         string
		mockResponse   string
		mockStatusCode int
		wantError      bool
		wantSymbol     string
	}{
		{
			name:           "successful stock request",
			symbol:         "DDOG",
			mockResponse:   testutils.YahooFinanceStockResponse,
			mockStatusCode: 200,
			wantError:      false,
			wantSymbol:     "DDOG",
		},
		{
			name:      "empty symbol",
			symbol:    "",
			wantError: true,
		},
		{
			name:      "invalid symbol",
			symbol:    "INVALID123",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := testutils.NewMockHTTPClient()
			service := NewService(mockClient)

			// Mock successful API response for valid symbols
			if !tt.wantError && tt.symbol != "" {
				expectedURL := "https://query1.finance.yahoo.com/v7/finance/quote?symbols=" + strings.ToUpper(tt.symbol)
				mockClient.AddResponse(expectedURL, tt.mockStatusCode, tt.mockResponse)
			}

			result, err := service.GetCurrentPrice(tt.symbol)

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
		})
	}
}

// TestService_FlakyRandomTest is a flaky test by design that fails roughly 50% of the time
func TestService_FlakyRandomTest(t *testing.T) {
	mockClient := testutils.NewMockHTTPClient()
	service := NewService(mockClient)

	// Mock successful API response
	expectedURL := "https://query1.finance.yahoo.com/v7/finance/quote?symbols=DDOG"
	mockClient.AddResponse(expectedURL, 200, testutils.YahooFinanceStockResponse)

	// Use current nanosecond time as seed for randomness
	rand.Seed(time.Now().UnixNano())
	randomValue := rand.Intn(10) + 1 // Random number between 1 and 10

	// Make a normal API call (this should work fine)
	_, err := service.GetCurrentPrice("DDOG")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// Flaky assertion: fail randomly based on the random seed
	if randomValue <= 5 {
		t.Errorf("Random flaky failure: got unlucky number %d (≤5), test fails by design", randomValue)
	} else {
		t.Logf("Random success: got lucky number %d (>5), test passes", randomValue)
	}
}

// TestService_AnotherFlakyTest is another flaky test that uses random failure
func TestService_AnotherFlakyTest(t *testing.T) {
	mockClient := testutils.NewMockHTTPClient()
	service := NewService(mockClient)

	// Mock response
	expectedURL := "https://query1.finance.yahoo.com/v7/finance/quote?symbols=TEST"
	mockClient.AddResponse(expectedURL, 200, testutils.YahooFinanceStockResponse)

	// Use a different seed source for variety
	rand.Seed(time.Now().UnixNano() + int64(time.Now().Second()))
	randomValue := rand.Intn(10) + 1 // Random number between 1 and 10

	// Make the API call
	_, err := service.GetCurrentPrice("TEST")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// Flaky assertion: fail randomly based on the random seed
	if randomValue <= 5 {
		t.Errorf("Another random flaky failure: rolled %d (≤5), test fails by design", randomValue)
	} else {
		t.Logf("Another random success: rolled %d (>5), test passes", randomValue)
	}
}

func TestService_GetDatadogPrice(t *testing.T) {
	mockClient := testutils.NewMockHTTPClient()
	service := NewService(mockClient)

	expectedURL := "https://query1.finance.yahoo.com/v7/finance/quote?symbols=DDOG"
	mockClient.AddResponse(expectedURL, 200, testutils.YahooFinanceStockResponse)

	result, err := service.GetDatadogPrice()

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

func TestService_GetStockSummary(t *testing.T) {
	tests := []struct {
		name         string
		symbol       string
		mockResponse string
		wantContains []string
	}{
		{
			name:         "positive change",
			symbol:       "DDOG",
			mockResponse: testutils.YahooFinanceStockResponse,
			wantContains: []string{"DDOG", "125.67", "↗", "up", "Market Open"},
		},
		{
			name:         "market closed",
			symbol:       "DDOG",
			mockResponse: testutils.YahooFinanceMarketClosed,
			wantContains: []string{"DDOG", "125.67", "↘", "down", "Market Closed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := testutils.NewMockHTTPClient()
			service := NewService(mockClient)

			expectedURL := "https://query1.finance.yahoo.com/v7/finance/quote?symbols=" + tt.symbol
			mockClient.AddResponse(expectedURL, 200, tt.mockResponse)

			summary, err := service.GetStockSummary(tt.symbol)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(summary, want) {
					t.Errorf("Expected summary to contain '%s', got: %s", want, summary)
				}
			}
		})
	}
}

func TestService_GetDatadogSummary(t *testing.T) {
	mockClient := testutils.NewMockHTTPClient()
	service := NewService(mockClient)

	expectedURL := "https://query1.finance.yahoo.com/v7/finance/quote?symbols=DDOG"
	mockClient.AddResponse(expectedURL, 200, testutils.YahooFinanceStockResponse)

	summary, err := service.GetDatadogSummary()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expectedParts := []string{"Datadog", "DDOG", "125.67", "↗", "up"}
	for _, part := range expectedParts {
		if !strings.Contains(summary, part) {
			t.Errorf("Expected summary to contain '%s', got: %s", part, summary)
		}
	}
}

func TestService_IsMarketOpen(t *testing.T) {
	tests := []struct {
		name         string
		symbol       string
		mockResponse string
		wantOpen     bool
	}{
		{
			name:         "market open",
			symbol:       "DDOG",
			mockResponse: testutils.YahooFinanceStockResponse,
			wantOpen:     true,
		},
		{
			name:         "market closed",
			symbol:       "DDOG",
			mockResponse: testutils.YahooFinanceMarketClosed,
			wantOpen:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := testutils.NewMockHTTPClient()
			service := NewService(mockClient)

			expectedURL := "https://query1.finance.yahoo.com/v7/finance/quote?symbols=" + tt.symbol
			mockClient.AddResponse(expectedURL, 200, tt.mockResponse)

			isOpen, err := service.IsMarketOpen(tt.symbol)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if isOpen != tt.wantOpen {
				t.Errorf("Expected market open %v, got %v", tt.wantOpen, isOpen)
			}
		})
	}
}

func TestService_GetPriceChange(t *testing.T) {
	tests := []struct {
		name         string
		symbol       string
		mockResponse string
		wantContains []string
	}{
		{
			name:         "positive change",
			symbol:       "DDOG",
			mockResponse: testutils.YahooFinanceStockResponse,
			wantContains: []string{"+2.34", "(1.89%)"},
		},
		{
			name:         "negative change",
			symbol:       "DDOG",
			mockResponse: testutils.YahooFinanceMarketClosed,
			wantContains: []string{"-1.23", "(-0.97%)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := testutils.NewMockHTTPClient()
			service := NewService(mockClient)

			expectedURL := "https://query1.finance.yahoo.com/v7/finance/quote?symbols=" + tt.symbol
			mockClient.AddResponse(expectedURL, 200, tt.mockResponse)

			change, err := service.GetPriceChange(tt.symbol)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(change, want) {
					t.Errorf("Expected change to contain '%s', got: %s", want, change)
				}
			}
		})
	}
}

func TestService_ValidateAndNormalizeSymbol(t *testing.T) {
	tests := []struct {
		name      string
		symbol    string
		want      string
		wantError bool
	}{
		{
			name:      "valid symbol",
			symbol:    "DDOG",
			want:      "DDOG",
			wantError: false,
		},
		{
			name:      "empty symbol",
			symbol:    "",
			wantError: true,
		},
		{
			name:      "invalid symbol with numbers",
			symbol:    "DD0G",
			wantError: true,
		},
	}

	service := NewService(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ValidateAndNormalizeSymbol(tt.symbol)

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

			if result != tt.want {
				t.Errorf("Expected result %v, got %v", tt.want, result)
			}
		})
	}
}

func TestNewService(t *testing.T) {
	t.Run("with custom client", func(t *testing.T) {
		mockClient := testutils.NewMockHTTPClient()
		service := NewService(mockClient)
		if service == nil {
			t.Errorf("Expected service, but got nil")
		}
		if service.client == nil {
			t.Errorf("Expected client to be set")
		}
	})

	t.Run("with nil client", func(t *testing.T) {
		service := NewService(nil)
		if service == nil {
			t.Errorf("Expected service, but got nil")
		}
		if service.client == nil {
			t.Errorf("Expected default client to be set")
		}
	})
}

// Test helper functions that are specific to models
func TestStockResponse_Methods(t *testing.T) {
	t.Run("IsPositiveChange", func(t *testing.T) {
		tests := []struct {
			name   string
			change float64
			want   bool
		}{
			{"positive change", 1.5, true},
			{"negative change", -1.5, false},
			{"zero change", 0.0, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				stock := &models.StockResponse{Change: tt.change}
				if got := stock.IsPositiveChange(); got != tt.want {
					t.Errorf("IsPositiveChange() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("GetChangeDirection", func(t *testing.T) {
		tests := []struct {
			name   string
			change float64
			want   string
		}{
			{"positive change", 1.5, "up"},
			{"negative change", -1.5, "down"},
			{"zero change", 0.0, "neutral"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				stock := &models.StockResponse{Change: tt.change}
				if got := stock.GetChangeDirection(); got != tt.want {
					t.Errorf("GetChangeDirection() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}
