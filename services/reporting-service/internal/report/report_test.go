package report

import "testing"

func TestBuildStatement(t *testing.T) {
	items := []LineItem{
		{Date: "2026-01-01", Description: "deposit", Amount: 5000},
		{Date: "2026-01-02", Description: "trade fee", Amount: -10},
	}
	s := BuildStatement("acc-1", items, 1000)
	if s.ClosingBalance != 5990 { t.Fatalf("want 5990, got %f", s.ClosingBalance) }
}
