package risk

import "testing"

func TestEvaluateHighRisk(t *testing.T) {
	got := Evaluate(Request{
		ExposureValue: 50000,
		PortfolioNAV:  100000,
		VaR95:         0.05,
		MaxDrawdown:   0.2,
		Leverage:      2.5,
	})
	if got.Band != "high" || !got.TradingBlock {
		t.Fatalf("expected high risk block, got %+v", got)
	}
}

func TestEvaluateLowRisk(t *testing.T) {
	got := Evaluate(Request{
		ExposureValue: 10000,
		PortfolioNAV:  100000,
		VaR95:         0.01,
		MaxDrawdown:   0.05,
		Leverage:      1.2,
	})
	if got.Band != "low" || got.TradingBlock {
		t.Fatalf("expected low risk without block, got %+v", got)
	}
}
