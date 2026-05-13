package marketdata

import (
	"testing"
	"time"
)

func TestBuildSnapshotHappyPath(t *testing.T) {
	got := BuildSnapshot(Request{
		Symbol:       "aapl",
		Bid:          100,
		Ask:          100.1,
		Last:         100.05,
		TickTimeUnix: time.Now().UTC().Unix(),
	})

	if !got.Valid || got.Symbol != "AAPL" {
		t.Fatalf("unexpected snapshot: %+v", got)
	}
}

func TestBuildSnapshotInvalidQuote(t *testing.T) {
	got := BuildSnapshot(Request{Symbol: "AAPL", Bid: 101, Ask: 100})
	if got.Valid || got.Reason != "invalid_quote" {
		t.Fatalf("expected invalid quote, got %+v", got)
	}
}

