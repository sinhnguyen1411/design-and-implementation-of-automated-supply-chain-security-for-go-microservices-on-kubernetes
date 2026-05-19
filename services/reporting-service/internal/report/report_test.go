package report_test

import (
	"testing"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/reporting-service/internal/report"
)

func TestFIFO_basic(t *testing.T) {
	var l report.Ledger
	l.Acquire(report.TaxLot{Symbol: "AAPL", Qty: 10, CostBasis: 100, AcquiredAt: "2022-01-01"})
	l.Acquire(report.TaxLot{Symbol: "AAPL", Qty: 5, CostBasis: 150, AcquiredAt: "2023-06-01"})

	gains, err := l.Sell(report.Sale{Symbol: "AAPL", Qty: 10, Proceeds: 200, SoldAt: "2024-01-15"})
	if err != nil {
		t.Fatal(err)
	}
	if len(gains) != 1 {
		t.Fatalf("expected 1 gain record (first lot only), got %d", len(gains))
	}
	if gains[0].GainLoss != 1000 {
		t.Fatalf("expected gain 1000 ((200-100)*10), got %f", gains[0].GainLoss)
	}
}

func TestFIFO_splitLots(t *testing.T) {
	var l report.Ledger
	l.Acquire(report.TaxLot{Symbol: "MSFT", Qty: 5, CostBasis: 200, AcquiredAt: "2021-01-01"})
	l.Acquire(report.TaxLot{Symbol: "MSFT", Qty: 5, CostBasis: 300, AcquiredAt: "2022-01-01"})

	gains, err := l.Sell(report.Sale{Symbol: "MSFT", Qty: 8, Proceeds: 400, SoldAt: "2024-06-01"})
	if err != nil {
		t.Fatal(err)
	}
	if len(gains) != 2 {
		t.Fatalf("expected 2 gain records for split lots, got %d", len(gains))
	}
	total := gains[0].GainLoss + gains[1].GainLoss
	expected := (400-200)*5.0 + (400-300)*3.0
	if total != expected {
		t.Fatalf("expected total gain %.0f, got %.0f", expected, total)
	}
}

func TestFIFO_shortVsLong(t *testing.T) {
	var l report.Ledger
	l.Acquire(report.TaxLot{Symbol: "GOOG", Qty: 1, CostBasis: 100, AcquiredAt: "2021-01-01"})
	l.Acquire(report.TaxLot{Symbol: "GOOG", Qty: 1, CostBasis: 200, AcquiredAt: "2023-12-01"})

	l.Sell(report.Sale{Symbol: "GOOG", Qty: 2, Proceeds: 300, SoldAt: "2024-06-01"})
	rep := l.Report()
	if rep.LongTerm == 0 {
		t.Fatal("expected non-zero long-term gain")
	}
	if rep.ShortTerm == 0 {
		t.Fatal("expected non-zero short-term gain")
	}
}

func TestFIFO_insufficient(t *testing.T) {
	var l report.Ledger
	l.Acquire(report.TaxLot{Symbol: "TSLA", Qty: 1, CostBasis: 100, AcquiredAt: "2023-01-01"})
	_, err := l.Sell(report.Sale{Symbol: "TSLA", Qty: 10, Proceeds: 200, SoldAt: "2024-01-01"})
	if err == nil {
		t.Fatal("expected insufficient shares error")
	}
}

func TestFIFO_report_totals(t *testing.T) {
	var l report.Ledger
	l.Acquire(report.TaxLot{Symbol: "NVDA", Qty: 5, CostBasis: 100, AcquiredAt: "2022-01-01"})
	l.Sell(report.Sale{Symbol: "NVDA", Qty: 5, Proceeds: 200, SoldAt: "2024-01-15"})

	rep := l.Report()
	if rep.TotalGain != 500 {
		t.Fatalf("expected total gain 500, got %f", rep.TotalGain)
	}
	if len(rep.Gains) != 1 {
		t.Fatalf("expected 1 gain record, got %d", len(rep.Gains))
	}
}
