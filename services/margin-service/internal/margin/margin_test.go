package margin

import "testing"

func TestCalculate(t *testing.T) {
	req := Req{Symbol: "AAPL", Qty: 100, Price: 150, MarginRatio: 0.25}
	r := Calculate(req, 10000)
	if r.Notional != 15000 { t.Fatalf("want 15000, got %f", r.Notional) }
	if r.InitialMargin != 3750 { t.Fatalf("want 3750, got %f", r.InitialMargin) }
	if r.ExcessLiquidity != 6250 { t.Fatalf("want 6250, got %f", r.ExcessLiquidity) }
}
