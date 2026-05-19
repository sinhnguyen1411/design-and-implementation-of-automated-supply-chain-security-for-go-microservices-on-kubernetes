package search_test

import (
	"testing"
	"time"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/search-service/internal/search"
)

var nyseHours = search.TradingHours{
	Timezone: "America/New_York",
	Open:     "09:30",
	Close:    "16:00",
	OpenDays: []string{"Mon", "Tue", "Wed", "Thu", "Fri"},
}

func makeAAPL() search.Symbol {
	return search.Symbol{
		Ticker: "AAPL", ISIN: "US0378331005", Name: "Apple Inc.",
		Exchange: "NASDAQ", Sector: "Technology", Currency: "USD",
		TradingHours: nyseHours,
	}
}

func TestAdd_validISIN(t *testing.T) {
	r := search.NewRegistry()
	if err := r.Add(makeAAPL()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAdd_invalidISIN(t *testing.T) {
	r := search.NewRegistry()
	s := makeAAPL()
	s.ISIN = "INVALID"
	if err := r.Add(s); err == nil {
		t.Fatal("expected error for invalid ISIN")
	}
}

func TestAdd_duplicate(t *testing.T) {
	r := search.NewRegistry()
	r.Add(makeAAPL())
	if err := r.Add(makeAAPL()); err == nil {
		t.Fatal("expected error for duplicate ISIN")
	}
}

func TestByTicker_caseInsensitive(t *testing.T) {
	r := search.NewRegistry()
	r.Add(makeAAPL())
	s, ok := r.ByTicker("aapl")
	if !ok {
		t.Fatal("ticker not found")
	}
	if s.ISIN != "US0378331005" {
		t.Fatalf("unexpected ISIN: %s", s.ISIN)
	}
}

func TestByISIN(t *testing.T) {
	r := search.NewRegistry()
	r.Add(makeAAPL())
	s, ok := r.ByISIN("US0378331005")
	if !ok {
		t.Fatal("ISIN not found")
	}
	if s.Ticker != "AAPL" {
		t.Fatalf("unexpected ticker: %s", s.Ticker)
	}
}

func TestSearch(t *testing.T) {
	r := search.NewRegistry()
	r.Add(makeAAPL())
	r.Add(search.Symbol{
		Ticker: "MSFT", ISIN: "US5949181045", Name: "Microsoft Corporation",
		Exchange: "NASDAQ", Sector: "Technology", Currency: "USD",
		TradingHours: nyseHours,
	})
	results := r.Search("apple")
	if len(results) != 1 || results[0].Ticker != "AAPL" {
		t.Fatalf("unexpected search results: %+v", results)
	}
}

func TestIsMarketOpen_closed(t *testing.T) {
	r := search.NewRegistry()
	r.Add(makeAAPL())
	// Monday 2024-01-15 12:00 UTC = 07:00 EST (before open)
	closed := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	open, err := r.IsMarketOpen("AAPL", closed)
	if err != nil {
		t.Fatal(err)
	}
	if open {
		t.Fatal("expected market closed at 07:00 EST")
	}
}

func TestIsMarketOpen_open(t *testing.T) {
	r := search.NewRegistry()
	r.Add(makeAAPL())
	// Monday 2024-01-15 15:00 UTC = 10:00 EST (during trading)
	during := time.Date(2024, 1, 15, 15, 0, 0, 0, time.UTC)
	open, err := r.IsMarketOpen("AAPL", during)
	if err != nil {
		t.Fatal(err)
	}
	if !open {
		t.Fatal("expected market open at 10:00 EST")
	}
}

func TestIsMarketOpen_weekend(t *testing.T) {
	r := search.NewRegistry()
	r.Add(makeAAPL())
	// Saturday 2024-01-13 15:00 UTC
	sat := time.Date(2024, 1, 13, 15, 0, 0, 0, time.UTC)
	open, err := r.IsMarketOpen("AAPL", sat)
	if err != nil {
		t.Fatal(err)
	}
	if open {
		t.Fatal("expected market closed on Saturday")
	}
}
