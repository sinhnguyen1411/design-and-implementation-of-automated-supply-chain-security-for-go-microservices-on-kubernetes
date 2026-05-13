package pricing

import "testing"

func TestBuildQuoteBuy(t *testing.T) {
	got := BuildQuote(Request{
		Symbol:         "msft",
		Side:           "buy",
		Quantity:       10,
		ReferencePrice: 50,
		FeeBps:         10,
		SpreadBps:      5,
	})
	if !got.Valid || got.Symbol != "MSFT" || got.GrossAmount <= 0 {
		t.Fatalf("unexpected quote: %+v", got)
	}
}

func TestBuildQuoteInvalidSide(t *testing.T) {
	got := BuildQuote(Request{
		Symbol:         "AAPL",
		Side:           "hold",
		Quantity:       1,
		ReferencePrice: 100,
	})
	if got.Valid || got.Reason != "invalid_side" {
		t.Fatalf("expected invalid side, got %+v", got)
	}
}

