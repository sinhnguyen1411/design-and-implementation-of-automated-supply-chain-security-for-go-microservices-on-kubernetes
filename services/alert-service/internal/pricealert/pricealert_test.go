package pricealert

import "testing"

func TestRegisterAndEvaluate(t *testing.T) {
	store := NewStore()
	store.Register("u1", "AAPL", 150, TypeAbove)
	store.Register("u1", "TSLA", 200, TypeBelow)

	events := store.EvaluateBatch(map[string]float64{"AAPL": 155, "TSLA": 190})
	if len(events) != 2 { t.Fatalf("want 2 events, got %d", len(events)) }
}

func TestDedup(t *testing.T) {
	store := NewStore()
	store.Register("u1", "AAPL", 150, TypeAbove)

	events1 := store.EvaluateBatch(map[string]float64{"AAPL": 160})
	if len(events1) != 1 { t.Fatal("first eval should fire") }

	events2 := store.EvaluateBatch(map[string]float64{"AAPL": 165})
	if len(events2) != 0 { t.Fatal("second eval should not fire (dedup)") }
}

func TestReset(t *testing.T) {
	store := NewStore()
	a := store.Register("u1", "MSFT", 300, TypeAbove)
	store.EvaluateBatch(map[string]float64{"MSFT": 310})

	ok := store.Reset(a.ID)
	if !ok { t.Fatal("reset should succeed") }

	events := store.EvaluateBatch(map[string]float64{"MSFT": 320})
	if len(events) != 1 { t.Fatal("after reset, should fire again") }
}

func TestPercentChange(t *testing.T) {
	store := NewStore()
	store.Register("u1", "NVDA", 5.0, TypePercentChange)
	events := store.EvaluateBatch(map[string]float64{"NVDA": 6.5})
	if len(events) != 1 { t.Fatal("6.5% change should trigger ±5% alert") }
}

func TestList(t *testing.T) {
	store := NewStore()
	store.Register("u1", "AAPL", 150, TypeAbove)
	store.Register("u1", "MSFT", 300, TypeBelow)
	store.Register("u2", "GOOGL", 100, TypeAbove)
	list := store.List("u1")
	if len(list) != 2 { t.Fatalf("want 2 alerts for u1, got %d", len(list)) }
}
