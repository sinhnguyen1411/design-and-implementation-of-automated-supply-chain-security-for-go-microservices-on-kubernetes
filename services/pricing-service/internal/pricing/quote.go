package pricing

import (
	"strings"

	"github.com/google/uuid"
)

type Request struct {
	Symbol         string  `json:"symbol"`
	Side           string  `json:"side"`
	Quantity       float64 `json:"quantity"`
	ReferencePrice float64 `json:"reference_price"`
	FeeBps         float64 `json:"fee_bps"`
	SpreadBps      float64 `json:"spread_bps"`
}

type Quote struct {
	QuoteID       string  `json:"quote_id"`
	Symbol        string  `json:"symbol"`
	Side          string  `json:"side"`
	AdjustedPrice float64 `json:"adjusted_price"`
	Notional      float64 `json:"notional"`
	Fee           float64 `json:"fee"`
	GrossAmount   float64 `json:"gross_amount"`
	Valid         bool    `json:"valid"`
	Reason        string  `json:"reason,omitempty"`
}

func BuildQuote(req Request) Quote {
	symbol := strings.ToUpper(strings.TrimSpace(req.Symbol))
	side := strings.ToLower(strings.TrimSpace(req.Side))
	out := Quote{
		QuoteID: uuid.NewString(),
		Symbol:  symbol,
		Side:    side,
	}

	if symbol == "" || req.Quantity <= 0 || req.ReferencePrice <= 0 {
		out.Valid = false
		out.Reason = "invalid_payload"
		return out
	}
	if side != "buy" && side != "sell" {
		out.Valid = false
		out.Reason = "invalid_side"
		return out
	}

	adjust := req.ReferencePrice * (req.SpreadBps / 10000.0)
	adjustedPrice := req.ReferencePrice
	if side == "buy" {
		adjustedPrice += adjust
	} else {
		adjustedPrice -= adjust
	}

	notional := adjustedPrice * req.Quantity
	fee := notional * (req.FeeBps / 10000.0)

	out.AdjustedPrice = adjustedPrice
	out.Notional = notional
	out.Fee = fee
	if side == "buy" {
		out.GrossAmount = notional + fee
	} else {
		out.GrossAmount = notional - fee
	}
	out.Valid = true
	return out
}

