package backtest

import "testing"

func TestRun(t *testing.T) {
	prices := []float64{100, 102, 101, 105, 103, 108}
	r := Run(prices, 10000)
	if r.Trades == 0 { t.Fatal("expected trades") }
	if r.FinalCapital <= 0 { t.Fatal("capital must be positive") }
}

func TestRun_insufficient_data(t *testing.T) {
	r := Run([]float64{100}, 5000)
	if r.FinalCapital != 5000 { t.Fatal("should return initial capital") }
}
