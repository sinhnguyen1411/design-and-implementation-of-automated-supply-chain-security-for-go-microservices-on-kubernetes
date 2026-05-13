package execution

import "testing"

func TestSimulatePartialFill(t *testing.T) {
	got := Simulate(Request{
		OrderQty:           100,
		AvailableLiquidity: 40,
		LimitPrice:         10,
		SlippageBps:        5,
	})
	if !got.Valid || got.Status != "partial_fill" || got.RemainingQty != 60 {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestSimulateInvalidOrder(t *testing.T) {
	got := Simulate(Request{OrderQty: 0, LimitPrice: 10})
	if got.Valid || got.Reason != "invalid_order" {
		t.Fatalf("expected invalid order, got %+v", got)
	}
}

