package stock

import (
	"fmt"
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

// TestService_ConcurrentRateLimiting is a flaky test by design that tests race conditions
// in rate limiting logic. This test may fail intermittently due to timing issues.
func TestService_ConcurrentRateLimiting(t *testing.T) {
	mockClient := testutils.NewMockHTTPClient()
	service := NewService(mockClient)

	// Mock successful API response
	expectedURL := "https://query1.finance.yahoo.com/v7/finance/quote?symbols=DDOG"
	mockClient.AddResponse(expectedURL, 200, testutils.YahooFinanceStockResponse)

	numGoroutines := 5
	results := make(chan error, numGoroutines)

	// Launch multiple goroutines to create race condition
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			// Add some random delay to increase chance of race condition
			time.Sleep(time.Duration(id*50) * time.Millisecond)

			start := time.Now()
			_, err := service.GetCurrentPrice("DDOG")
			duration := time.Since(start)

			// This assertion is flaky by design - it expects all requests
			// to complete within 3 seconds, but with rate limiting (2s between requests)
			// and multiple concurrent requests, this may fail randomly
			if duration > 3*time.Second {
				results <- fmt.Errorf("request %d took too long: %v", id, duration)
				return
			}

			results <- err
		}(i)
	}

	// Collect results
	var errors []error
	for i := 0; i < numGoroutines; i++ {
		if err := <-results; err != nil {
			errors = append(errors, err)
		}
	}

	// This test is designed to be flaky - sometimes it passes, sometimes it fails
	// depending on exact timing of goroutines and rate limiting behavior
	if len(errors) > 2 { // Allow some failures but not all
		t.Errorf("Too many concurrent requests failed (%d/%d): %v", len(errors), numGoroutines, errors)
	}
}

// TestService_TimingDependentValidation is another flaky test that depends on system timing
func TestService_TimingDependentValidation(t *testing.T) {
	mockClient := testutils.NewMockHTTPClient()
	service := NewService(mockClient)

	// Mock response
	expectedURL := "https://query1.finance.yahoo.com/v7/finance/quote?symbols=TEST"
	mockClient.AddResponse(expectedURL, 200, testutils.YahooFinanceStockResponse)

	// This test uses nanosecond timing which is inherently flaky
	start := time.Now()
	_, err := service.GetCurrentPrice("TEST")
	nanos := time.Since(start).Nanoseconds()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// More aggressive flaky assertion: use multiple criteria to increase failure chance
	// This will fail roughly 30% of the time based on nanosecond timing patterns
	lastDigit := nanos % 10
	secondLastDigit := (nanos / 10) % 10

	if lastDigit == 0 || lastDigit == 3 || lastDigit == 7 || secondLastDigit == 2 {
		t.Errorf("Request completed at an 'unlucky' nanosecond timing: %d ns (last digits: %d%d)", nanos, secondLastDigit, lastDigit)
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
