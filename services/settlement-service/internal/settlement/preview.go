package settlement

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Request struct {
	TradeDate   string  `json:"trade_date"`
	CycleDays   int     `json:"cycle_days"`
	Side        string  `json:"side"`
	GrossAmount float64 `json:"gross_amount"`
	Fee         float64 `json:"fee"`
	TaxRate     float64 `json:"tax_rate"`
}

type Preview struct {
	PreviewID       string  `json:"preview_id"`
	SettlementDate  string  `json:"settlement_date"`
	TaxAmount       float64 `json:"tax_amount"`
	NetCashMovement float64 `json:"net_cash_movement"`
	Status          string  `json:"status"`
	Valid           bool    `json:"valid"`
	Reason          string  `json:"reason,omitempty"`
}

func BuildPreview(req Request) Preview {
	side := strings.ToLower(strings.TrimSpace(req.Side))
	out := Preview{PreviewID: uuid.NewString()}

	tradeDate, err := time.Parse("2006-01-02", req.TradeDate)
	if err != nil || req.CycleDays < 0 || req.GrossAmount <= 0 {
		out.Valid = false
		out.Reason = "invalid_payload"
		return out
	}
	if side != "buy" && side != "sell" {
		out.Valid = false
		out.Reason = "invalid_side"
		return out
	}

	if req.TaxRate < 0 {
		req.TaxRate = 0
	}
	tax := req.GrossAmount * req.TaxRate
	net := req.GrossAmount - req.Fee - tax
	if side == "buy" {
		net = -(req.GrossAmount + req.Fee + tax)
	}

	settleDate := tradeDate.AddDate(0, 0, req.CycleDays)
	status := "pending"
	if !settleDate.After(time.Now().UTC()) {
		status = "due"
	}

	out.SettlementDate = settleDate.Format("2006-01-02")
	out.TaxAmount = tax
	out.NetCashMovement = net
	out.Status = status
	out.Valid = true
	return out
}

