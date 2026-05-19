package report

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TaxLot struct {
	LotID      string  `json:"lot_id"`
	Symbol     string  `json:"symbol"`
	Qty        float64 `json:"qty"`
	CostBasis  float64 `json:"cost_basis"`
	AcquiredAt string  `json:"acquired_at"`
}

type Sale struct {
	Symbol   string  `json:"symbol"`
	Qty      float64 `json:"qty"`
	Proceeds float64 `json:"proceeds"`
	SoldAt   string  `json:"sold_at"`
}

type RealizedGain struct {
	LotID       string  `json:"lot_id"`
	Symbol      string  `json:"symbol"`
	Qty         float64 `json:"qty"`
	CostBasis   float64 `json:"cost_basis"`
	Proceeds    float64 `json:"proceeds"`
	GainLoss    float64 `json:"gain_loss"`
	HoldingDays int     `json:"holding_days"`
	Term        string  `json:"term"`
}

type TaxReport struct {
	Gains     []RealizedGain `json:"gains"`
	TotalGain float64        `json:"total_gain"`
	ShortTerm float64        `json:"short_term"`
	LongTerm  float64        `json:"long_term"`
}

type Ledger struct {
	lots          []TaxLot
	realizedGains []RealizedGain
}

var (
	ErrInsufficientShares = errors.New("insufficient shares for sale")
	ErrInvalidQty         = errors.New("quantity must be positive")
)

func (l *Ledger) Acquire(lot TaxLot) {
	if lot.LotID == "" {
		lot.LotID = uuid.NewString()
	}
	l.lots = append(l.lots, lot)
}

func (l *Ledger) Sell(sale Sale) ([]RealizedGain, error) {
	if sale.Qty <= 0 {
		return nil, ErrInvalidQty
	}
	available := 0.0
	for _, lot := range l.lots {
		if lot.Symbol == sale.Symbol && lot.Qty > 0 {
			available += lot.Qty
		}
	}
	if available < sale.Qty {
		return nil, ErrInsufficientShares
	}
	soldAt, _ := time.Parse("2006-01-02", sale.SoldAt)
	remaining := sale.Qty
	var gains []RealizedGain
	for i := range l.lots {
		if remaining <= 0 {
			break
		}
		lot := &l.lots[i]
		if lot.Symbol != sale.Symbol || lot.Qty <= 0 {
			continue
		}
		qty := lot.Qty
		if qty > remaining {
			qty = remaining
		}
		lot.Qty -= qty
		remaining -= qty
		acquiredAt, _ := time.Parse("2006-01-02", lot.AcquiredAt)
		holdingDays := int(soldAt.Sub(acquiredAt).Hours() / 24)
		term := "short"
		if holdingDays > 365 {
			term = "long"
		}
		gains = append(gains, RealizedGain{
			LotID:       lot.LotID,
			Symbol:      sale.Symbol,
			Qty:         qty,
			CostBasis:   lot.CostBasis,
			Proceeds:    sale.Proceeds,
			GainLoss:    (sale.Proceeds - lot.CostBasis) * qty,
			HoldingDays: holdingDays,
			Term:        term,
		})
	}
	l.realizedGains = append(l.realizedGains, gains...)
	return gains, nil
}

func (l *Ledger) Report() TaxReport {
	total, short, long := 0.0, 0.0, 0.0
	for _, g := range l.realizedGains {
		total += g.GainLoss
		if g.Term == "long" {
			long += g.GainLoss
		} else {
			short += g.GainLoss
		}
	}
	return TaxReport{
		Gains:     l.realizedGains,
		TotalGain: total,
		ShortTerm: short,
		LongTerm:  long,
	}
}
