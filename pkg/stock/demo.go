package stock

import (
	"math/rand"
	"time"

	"github.com/JSGette/agent_summit_bazel_workshop/pkg/models"
)

// DemoStockData contains realistic demo data for stocks
var DemoStockData = map[string]struct {
	Name      string
	BasePrice float64
	Currency  string
	MarketCap int64
}{
	"DDOG": {
		Name:      "Datadog, Inc.",
		BasePrice: 125.50,
		Currency:  "USD",
		MarketCap: 40000000000,
	},
	"AAPL": {
		Name:      "Apple Inc.",
		BasePrice: 175.25,
		Currency:  "USD",
		MarketCap: 2800000000000,
	},
	"GOOGL": {
		Name:      "Alphabet Inc.",
		BasePrice: 142.75,
		Currency:  "USD",
		MarketCap: 1800000000000,
	},
	"MSFT": {
		Name:      "Microsoft Corporation",
		BasePrice: 415.50,
		Currency:  "USD",
		MarketCap: 3100000000000,
	},
	"TSLA": {
		Name:      "Tesla, Inc.",
		BasePrice: 248.75,
		Currency:  "USD",
		MarketCap: 790000000000,
	},
}

// generateDemoStockResponse creates a realistic stock response with simulated price movements
func generateDemoStockResponse(symbol string) (*models.StockResponse, error) {
	data, exists := DemoStockData[symbol]
	if !exists {
		return nil, models.NewAPIError("Demo Stock", "Stock symbol not found in demo data", 404)
	}

	// Create a deterministic but varying price based on current time
	now := time.Now()
	seed := now.Hour()*60 + now.Minute() // Changes every minute
	r := rand.New(rand.NewSource(int64(seed + len(symbol))))

	// Generate price variation (-5% to +5%)
	variation := (r.Float64() - 0.5) * 0.1 // -0.05 to +0.05
	currentPrice := data.BasePrice * (1 + variation)

	// Calculate change from "yesterday"
	yesterdayVariation := (r.Float64() - 0.5) * 0.08 // Slightly smaller range for yesterday
	yesterdayPrice := data.BasePrice * (1 + yesterdayVariation)
	change := currentPrice - yesterdayPrice
	changePercent := (change / yesterdayPrice) * 100

	// Generate volume (random but reasonable)
	volume := int64(500000 + r.Intn(2000000)) // 500K to 2.5M shares

	// Determine market state based on time (simplified)
	var marketState models.MarketState
	hour := now.Hour()
	if hour >= 9 && hour < 16 { // 9 AM to 4 PM (simplified US market hours)
		marketState = models.MarketStateRegular
	} else if hour >= 4 && hour < 9 { // Pre-market
		marketState = models.MarketStatePremarket
	} else if hour >= 16 && hour < 20 { // After hours
		marketState = models.MarketStatePostmarket
	} else {
		marketState = models.MarketStateClosed
	}

	return &models.StockResponse{
		Symbol:        symbol,
		CompanyName:   data.Name,
		Price:         currentPrice,
		Change:        change,
		ChangePercent: changePercent,
		PreviousClose: yesterdayPrice,
		Volume:        volume,
		MarketCap:     data.MarketCap,
		MarketState:   marketState,
		Currency:      data.Currency,
		Metadata: models.ResponseMetadata{
			Timestamp: now,
			Source:    "Demo Mode (Simulated Data)",
		},
	}, nil
}

// GetDemoStock returns demo stock data for the given symbol
func GetDemoStock(symbol string) (*models.StockResponse, error) {
	return generateDemoStockResponse(symbol)
}
