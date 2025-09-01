package testutils

import (
	"bytes"
	"io"
	"net/http"
)

// MockHTTPClient is a mock implementation of HTTPClient for testing
type MockHTTPClient struct {
	Responses map[string]*http.Response
	Errors    map[string]error
	CallCount map[string]int
}

// NewMockHTTPClient creates a new mock HTTP client
func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{
		Responses: make(map[string]*http.Response),
		Errors:    make(map[string]error),
		CallCount: make(map[string]int),
	}
}

// Get implements the HTTPClient interface
func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	m.CallCount[url]++

	if err, exists := m.Errors[url]; exists {
		return nil, err
	}

	if resp, exists := m.Responses[url]; exists {
		return resp, nil
	}

	// Default response for unmocked URLs
	return &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"error": "Not found"}`))),
	}, nil
}

// AddResponse adds a mock response for a given URL
func (m *MockHTTPClient) AddResponse(url string, statusCode int, body string) {
	m.Responses[url] = &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Header:     make(http.Header),
	}
}

// AddError adds a mock error for a given URL
func (m *MockHTTPClient) AddError(url string, err error) {
	m.Errors[url] = err
}

// GetCallCount returns the number of times a URL was called
func (m *MockHTTPClient) GetCallCount(url string) int {
	return m.CallCount[url]
}

// Reset clears all mock data
func (m *MockHTTPClient) Reset() {
	m.Responses = make(map[string]*http.Response)
	m.Errors = make(map[string]error)
	m.CallCount = make(map[string]int)
}
