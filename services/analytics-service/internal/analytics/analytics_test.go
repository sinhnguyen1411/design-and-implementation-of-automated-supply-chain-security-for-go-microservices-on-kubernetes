package analytics_test

import (
	"math"
	"testing"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/analytics-service/internal/analytics"
)

func TestComputeMetrics_empty(t *testing.T) {
	m := analytics.ComputeMetrics(nil)
	if m.Sharpe != 0 || m.CAGR != 0 {
		t.Fatalf("expected zero metrics for empty input, got %+v", m)
	}
}

func TestComputeMetrics_uniformPositive(t *testing.T) {
	var returns []analytics.DailyReturn
	for i := 0; i < 252; i++ {
		returns = append(returns, analytics.DailyReturn{Date: "2024", Return: 0.01})
	}
	m := analytics.ComputeMetrics(returns)
	if m.TotalReturn < 10 {
		t.Fatalf("expected high total return (1.01^252-1), got %.4f", m.TotalReturn)
	}
	// Drawdown is non-negative by construction; on arm64 (macOS) FP contraction
	// yields ~1e-18 instead of exactly 0, so assert with an epsilon tolerance
	// rather than strict equality to keep the cross-OS matrix green.
	if m.MaxDrawdown > 1e-9 {
		t.Fatalf("expected ~zero drawdown for monotone positive returns, got %.4e", m.MaxDrawdown)
	}
}

func TestComputeMetrics_drawdown(t *testing.T) {
	returns := []analytics.DailyReturn{
		{Return: 0.1},
		{Return: 0.1},
		{Return: -0.5},
		{Return: 0.1},
	}
	m := analytics.ComputeMetrics(returns)
	if m.MaxDrawdown <= 0 {
		t.Fatalf("expected positive max drawdown, got %.4f", m.MaxDrawdown)
	}
	if math.IsNaN(m.CAGR) {
		t.Fatal("CAGR is NaN")
	}
}

func TestComputeMetrics_sharpePositive(t *testing.T) {
	var returns []analytics.DailyReturn
	for i := 0; i < 100; i++ {
		returns = append(returns, analytics.DailyReturn{Return: 0.005})
	}
	for i := 0; i < 100; i++ {
		returns = append(returns, analytics.DailyReturn{Return: -0.001})
	}
	m := analytics.ComputeMetrics(returns)
	if m.Sharpe <= 0 {
		t.Fatalf("expected positive Sharpe ratio, got %.4f", m.Sharpe)
	}
}

func TestComputeMetrics_volatility(t *testing.T) {
	var returns []analytics.DailyReturn
	for i := 0; i < 50; i++ {
		returns = append(returns, analytics.DailyReturn{Return: 0.01})
		returns = append(returns, analytics.DailyReturn{Return: -0.01})
	}
	m := analytics.ComputeMetrics(returns)
	if m.Volatility <= 0 {
		t.Fatalf("expected positive volatility, got %.4f", m.Volatility)
	}
}
