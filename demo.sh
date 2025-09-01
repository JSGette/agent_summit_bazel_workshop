#!/bin/bash

# Weather & Stock API Demo Script
# This script demonstrates the API endpoints with real examples

BASE_URL="http://localhost:3000"
echo "🌤️  Weather & Stock API Demo"
echo "=================================="
echo ""

# Function to make API calls with nice formatting
api_call() {
    local url="$1"
    local description="$2"
    
    echo "📡 $description"
    echo "   GET $url"
    echo "   Response:"
    echo "   --------"
    curl -s "$url" | jq '.' || echo "   (Failed to connect - is the server running?)"
    echo ""
    echo ""
}

# Check if server is running
echo "🔍 Checking if server is running..."
if ! curl -s "$BASE_URL/health" > /dev/null; then
    echo "❌ Server is not running on $BASE_URL"
    echo "   Please start the server first:"
    echo "   ./bin/weather-stock-api"
    echo ""
    exit 1
fi
echo "✅ Server is running!"
echo ""

# Demo API endpoints
echo "🎯 API Demonstration"
echo "====================="
echo ""

# Health check
api_call "$BASE_URL/health" "Health Check"

# API info
api_call "$BASE_URL/" "API Information"

# Weather examples
api_call "$BASE_URL/weather?city=Stuttgart" "Weather for Stuttgart"
api_call "$BASE_URL/weather?city=Berlin" "Weather for Berlin"
api_call "$BASE_URL/weather/summary?city=New York" "Weather Summary for New York"

# Stock examples
api_call "$BASE_URL/stock/datadog" "Datadog Stock Price"
api_call "$BASE_URL/stock?symbol=AAPL" "Apple Stock Price"
api_call "$BASE_URL/stock/summary?symbol=DDOG" "Datadog Stock Summary"

echo "🎉 Demo completed!"
echo ""
echo "💡 Try these URLs in your browser:"
echo "   $BASE_URL/"
echo "   $BASE_URL/health"
echo "   $BASE_URL/weather?city=Stuttgart"
echo "   $BASE_URL/stock/datadog"
echo ""
