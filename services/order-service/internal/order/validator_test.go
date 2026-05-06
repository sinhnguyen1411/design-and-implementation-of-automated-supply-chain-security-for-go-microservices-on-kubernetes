package order

import "testing"

func TestValidateBuyInsufficientCash(t *testing.T) {
	got := Validate(Request{
		Symbol:      "AAPL",
		Side:        "buy",
		Qty:         10,
		Price:       100,
		CashBalance: 900,
	})
	if got.Accepted || got.Reason != "insufficient_cash" {
		t.Fatalf("unexpected decision: %+v", got)
	}
}

func TestValidateSellAccepted(t *testing.T) {
	got := Validate(Request{
		Symbol:     "AAPL",
		Side:       "sell",
		Qty:        2,
		Price:      100,
		HoldingQty: 10,
	})
	if !got.Accepted || got.Reason != "ok" {
		t.Fatalf("unexpected decision: %+v", got)
	}
}
