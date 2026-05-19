package datafeed

import "testing"

func TestAggregateToBar(t *testing.T) {
	ticks := []Tick{{Symbol: "AAPL", Price: 150, Volume: 100}, {Symbol: "AAPL", Price: 155, Volume: 200}, {Symbol: "AAPL", Price: 148, Volume: 150}}
	b := AggregateToBar(ticks)
	if b.High != 155 { t.Fatalf("want high=155, got %f", b.High) }
	if b.Low != 148 { t.Fatalf("want low=148, got %f", b.Low) }
	if b.Volume != 450 { t.Fatalf("want vol=450, got %f", b.Volume) }
}

func TestAggregateToBar_empty(t *testing.T) {
	b := AggregateToBar(nil)
	if b.Open != 0 { t.Fatal("expected zero bar") }
}
