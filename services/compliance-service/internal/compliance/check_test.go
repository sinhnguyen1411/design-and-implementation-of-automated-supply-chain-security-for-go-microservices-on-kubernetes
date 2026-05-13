package compliance

import "testing"

func TestCheckAllowed(t *testing.T) {
	got := Check(Request{
		Symbol:            "AAPL",
		Side:              "buy",
		Qty:               10,
		Price:             100,
		CurrentPosition:   20,
		MaxPosition:       100,
		MaxNotional:       5000,
		RestrictedSymbols: []string{"TSLA"},
	})
	if !got.Valid || !got.Allowed {
		t.Fatalf("expected allowed, got %+v", got)
	}
}

func TestCheckRestrictedSymbol(t *testing.T) {
	got := Check(Request{
		Symbol:            "TSLA",
		Side:              "buy",
		Qty:               1,
		Price:             100,
		RestrictedSymbols: []string{"TSLA"},
	})
	if !got.Valid || got.Allowed {
		t.Fatalf("expected blocked symbol, got %+v", got)
	}
}

