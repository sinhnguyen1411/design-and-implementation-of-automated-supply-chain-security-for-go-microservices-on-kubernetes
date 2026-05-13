package settlement

import "testing"

func TestBuildPreviewSell(t *testing.T) {
	got := BuildPreview(Request{
		TradeDate:   "2026-05-13",
		CycleDays:   2,
		Side:        "sell",
		GrossAmount: 1000,
		Fee:         5,
		TaxRate:     0.01,
	})
	if !got.Valid || got.NetCashMovement <= 0 {
		t.Fatalf("unexpected preview: %+v", got)
	}
}

func TestBuildPreviewInvalidDate(t *testing.T) {
	got := BuildPreview(Request{TradeDate: "13-05-2026", CycleDays: 2, Side: "buy", GrossAmount: 100})
	if got.Valid || got.Reason != "invalid_payload" {
		t.Fatalf("expected invalid payload, got %+v", got)
	}
}

