package analytics

import "testing"

func TestComputePnL(t *testing.T) {
	trades := []Trade{
		{Symbol: "AAPL", Qty: 10, BuyPrice: 100, SellPrice: 120},
		{Symbol: "MSFT", Qty: 5, BuyPrice: 200, SellPrice: 190},
	}
	r := ComputePnL(trades, 0.001)
	if r.Trades != 2 { t.Fatalf("want 2, got %d", r.Trades) }
	if r.GrossPnL != 150 { t.Fatalf("want 150, got %f", r.GrossPnL) }
}

func TestComputePnL_empty(t *testing.T) {
	r := ComputePnL(nil, 0.001)
	if r.Trades != 0 { t.Fatal("expected zero result") }
}
