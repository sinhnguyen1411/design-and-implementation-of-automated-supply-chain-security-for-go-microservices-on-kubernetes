package alert

import "testing"

func TestEvaluate_fired(t *testing.T) {
	r := Evaluate(Alert{SymbolID: "AAPL", Threshold: 150, Type: AlertAbove}, 155)
	if !r.Fired { t.Fatal("should have fired") }
}

func TestEvaluate_not_fired(t *testing.T) {
	r := Evaluate(Alert{SymbolID: "AAPL", Threshold: 150, Type: AlertAbove}, 140)
	if r.Fired { t.Fatal("should not fire") }
}

func TestEvaluate_below(t *testing.T) {
	r := Evaluate(Alert{SymbolID: "TSLA", Threshold: 200, Type: AlertBelow}, 180)
	if !r.Fired { t.Fatal("should have fired below") }
}
