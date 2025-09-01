package models

import "time"

// MarketState represents the current state of the stock market
type MarketState string

const (
	MarketStateRegular    MarketState = "REGULAR"
	MarketStatePremarket  MarketState = "PRE"
	MarketStatePostmarket MarketState = "POST"
	MarketStateClosed     MarketState = "CLOSED"
)

// StockResponse represents the standardized stock response
type StockResponse struct {
	Symbol        string           `json:"symbol"`
	CompanyName   string           `json:"company_name"`
	Price         float64          `json:"price"`
	Change        float64          `json:"change"`
	ChangePercent float64          `json:"change_percent"`
	PreviousClose float64          `json:"previous_close"`
	Volume        int64            `json:"volume"`
	MarketCap     int64            `json:"market_cap,omitempty"`
	MarketState   MarketState      `json:"market_state"`
	Currency      string           `json:"currency"`
	Metadata      ResponseMetadata `json:"metadata"`
}

// YahooFinanceResponse represents the raw response from Yahoo Finance API
type YahooFinanceResponse struct {
	QuoteResponse struct {
		Result []struct {
			Symbol                     string  `json:"symbol"`
			ShortName                  string  `json:"shortName"`
			LongName                   string  `json:"longName"`
			RegularMarketPrice         float64 `json:"regularMarketPrice"`
			RegularMarketChange        float64 `json:"regularMarketChange"`
			RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
			RegularMarketPreviousClose float64 `json:"regularMarketPreviousClose"`
			RegularMarketVolume        int64   `json:"regularMarketVolume"`
			MarketCap                  int64   `json:"marketCap"`
			Currency                   string  `json:"currency"`
			MarketState                string  `json:"marketState"`
			RegularMarketTime          int64   `json:"regularMarketTime"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"quoteResponse"`
}

// ConvertYahooFinanceResponse converts Yahoo Finance API response to our standard format
func ConvertYahooFinanceResponse(response *YahooFinanceResponse) (*StockResponse, error) {
	if len(response.QuoteResponse.Result) == 0 {
		return nil, NewAPIError("Yahoo Finance", "No stock data found", 404)
	}

	result := response.QuoteResponse.Result[0]

	// Convert market state
	var marketState MarketState
	switch result.MarketState {
	case "REGULAR":
		marketState = MarketStateRegular
	case "PRE":
		marketState = MarketStatePremarket
	case "POST":
		marketState = MarketStatePostmarket
	case "CLOSED":
		marketState = MarketStateClosed
	default:
		marketState = MarketStateClosed
	}

	// Use long name if available, otherwise short name
	companyName := result.LongName
	if companyName == "" {
		companyName = result.ShortName
	}

	// Convert Unix timestamp to time
	timestamp := time.Unix(result.RegularMarketTime, 0)

	return &StockResponse{
		Symbol:        result.Symbol,
		CompanyName:   companyName,
		Price:         result.RegularMarketPrice,
		Change:        result.RegularMarketChange,
		ChangePercent: result.RegularMarketChangePercent,
		PreviousClose: result.RegularMarketPreviousClose,
		Volume:        result.RegularMarketVolume,
		MarketCap:     result.MarketCap,
		MarketState:   marketState,
		Currency:      result.Currency,
		Metadata: ResponseMetadata{
			Timestamp: timestamp,
			Source:    "Yahoo Finance",
		},
	}, nil
}

// IsPositiveChange returns true if the stock price change is positive
func (s *StockResponse) IsPositiveChange() bool {
	return s.Change > 0
}

// GetChangeDirection returns "up", "down", or "neutral" based on price change
func (s *StockResponse) GetChangeDirection() string {
	if s.Change > 0 {
		return "up"
	} else if s.Change < 0 {
		return "down"
	}
	return "neutral"
}
