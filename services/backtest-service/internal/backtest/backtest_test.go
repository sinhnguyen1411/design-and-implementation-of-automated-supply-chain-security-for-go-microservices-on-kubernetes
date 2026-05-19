package backtest_test

import (
	"testing"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/backtest-service/internal/backtest"
)

func TestRunBacktest_noPrices(t *testing.T) {
	strategy := &backtest.SMAStrategy{Short: 5, Long: 20}
	result := backtest.RunBacktest(strategy, nil, 10000)
	if result.FinalCapital != 10000 {
		t.Fatalf("expected 10000 capital with no prices, got %.2f", result.FinalCapital)
	}
}

func TestBuyAndHold_return(t *testing.T) {
	prices := []float64{100, 110, 120, 130, 140}
	result := backtest.BuyAndHold(prices, 10000)
	if result.TotalReturn < 0.39 || result.TotalReturn > 0.41 {
		t.Fatalf("expected total return ~0.40, got %.4f", result.TotalReturn)
	}
}

func TestBuyAndHold_maxDrawdown(t *testing.T) {
	prices := []float64{100, 120, 80, 110}
	result := backtest.BuyAndHold(prices, 10000)
	if result.MaxDrawdown <= 0 {
		t.Fatalf("expected positive max drawdown, got %.4f", result.MaxDrawdown)
	}
}

func TestSMAStrategy_name(t *testing.T) {
	strategy := &backtest.SMAStrategy{Short: 3, Long: 5}
	prices := []float64{10, 10, 10, 10, 10, 11, 12, 13, 14, 15}
	result := backtest.RunBacktest(strategy, prices, 10000)
	if result.StrategyName != "sma" {
		t.Fatalf("unexpected strategy name: %s", result.StrategyName)
	}
}

func TestMomentumStrategy_name(t *testing.T) {
	strategy := &backtest.MomentumStrategy{Window: 3, Threshold: 0.05}
	prices := []float64{100, 100, 100, 115, 130, 145}
	result := backtest.RunBacktest(strategy, prices, 10000)
	if result.StrategyName != "momentum" {
		t.Fatalf("unexpected strategy name: %s", result.StrategyName)
	}
}

func TestWalkForward_folds(t *testing.T) {
	strategy := &backtest.MomentumStrategy{Window: 2, Threshold: 0.02}
	var prices []float64
	for i := 0; i < 100; i++ {
		prices = append(prices, float64(100+i))
	}
	result := backtest.WalkForward(strategy, prices, 10000, 4)
	if result.Folds != 4 {
		t.Fatalf("expected 4 folds, got %d", result.Folds)
	}
	if len(result.OutOfSample) != 4 {
		t.Fatalf("expected 4 OOS results, got %d", len(result.OutOfSample))
	}
}

func TestWalkForward_benchmark(t *testing.T) {
	strategy := &backtest.SMAStrategy{Short: 2, Long: 5}
	var prices []float64
	for i := 0; i < 60; i++ {
		prices = append(prices, float64(100+i*2))
	}
	result := backtest.WalkForward(strategy, prices, 10000, 3)
	if result.Benchmark.TotalReturn <= 0 {
		t.Fatalf("expected positive benchmark return for rising prices, got %.4f", result.Benchmark.TotalReturn)
	}
}
