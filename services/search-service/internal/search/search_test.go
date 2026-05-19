package search

import "testing"

var catalog = []Symbol{
	{Ticker: "AAPL", Name: "Apple Inc", Exchange: "NASDAQ", Sector: "Technology"},
	{Ticker: "MSFT", Name: "Microsoft Corp", Exchange: "NASDAQ", Sector: "Technology"},
	{Ticker: "JPM", Name: "JPMorgan Chase", Exchange: "NYSE", Sector: "Finance"},
}

func TestQuery(t *testing.T) {
	r := Query(catalog, "app")
	if len(r) != 1 { t.Fatalf("want 1, got %d", len(r)) }
}

func TestQuery_empty(t *testing.T) {
	r := Query(catalog, "")
	if len(r) != len(catalog) { t.Fatal("empty query should return all") }
}
