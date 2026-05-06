package portfolio

import "testing"

func TestBuildSummary(t *testing.T) {
	positions := []Position{
		{Symbol: "AAPL", Qty: 10, Price: 200},
		{Symbol: "MSFT", Qty: 5, Price: 100},
	}

	got := BuildSummary(500, positions)

	if got.MarketValue != 2500 {
		t.Fatalf("market value mismatch: got %v", got.MarketValue)
	}
	if got.TotalEquity != 3000 {
		t.Fatalf("total equity mismatch: got %v", got.TotalEquity)
	}
	if got.SectorDiversifyHint != "high_concentration" {
		t.Fatalf("hint mismatch: got %s", got.SectorDiversifyHint)
	}
}
