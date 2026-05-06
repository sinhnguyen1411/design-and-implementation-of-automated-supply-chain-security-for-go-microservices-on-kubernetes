package order

import "strings"

type Request struct {
	Symbol       string  `json:"symbol"`
	Side         string  `json:"side"`
	Qty          float64 `json:"qty"`
	Price        float64 `json:"price"`
	CashBalance  float64 `json:"cash_balance"`
	HoldingQty   float64 `json:"holding_qty"`
	MaxOrderNotl float64 `json:"max_order_notional"`
}

type Decision struct {
	Accepted bool    `json:"accepted"`
	Reason   string  `json:"reason"`
	Notional float64 `json:"notional"`
	Fee      float64 `json:"fee"`
}

func Validate(req Request) Decision {
	if req.Symbol == "" || req.Qty <= 0 || req.Price <= 0 {
		return Decision{Accepted: false, Reason: "invalid_payload"}
	}

	side := strings.ToLower(strings.TrimSpace(req.Side))
	notional := req.Qty * req.Price
	fee := notional * 0.001

	if req.MaxOrderNotl > 0 && notional > req.MaxOrderNotl {
		return Decision{Accepted: false, Reason: "max_notional_exceeded", Notional: notional, Fee: fee}
	}

	switch side {
	case "buy":
		if req.CashBalance < notional+fee {
			return Decision{Accepted: false, Reason: "insufficient_cash", Notional: notional, Fee: fee}
		}
	case "sell":
		if req.HoldingQty < req.Qty {
			return Decision{Accepted: false, Reason: "insufficient_holding", Notional: notional, Fee: fee}
		}
	default:
		return Decision{Accepted: false, Reason: "invalid_side", Notional: notional, Fee: fee}
	}

	return Decision{Accepted: true, Reason: "ok", Notional: notional, Fee: fee}
}
