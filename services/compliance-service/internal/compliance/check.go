package compliance

import (
	"math"
	"slices"
	"strings"

	"github.com/google/uuid"
)

type Request struct {
	Symbol            string   `json:"symbol"`
	Side              string   `json:"side"`
	Qty               float64  `json:"qty"`
	Price             float64  `json:"price"`
	CurrentPosition   float64  `json:"current_position"`
	MaxPosition       float64  `json:"max_position"`
	MaxNotional       float64  `json:"max_notional"`
	RestrictedSymbols []string `json:"restricted_symbols"`
}

type Result struct {
	CheckID     string   `json:"check_id"`
	Allowed     bool     `json:"allowed"`
	Notional    float64  `json:"notional"`
	Violations  []string `json:"violations"`
	Valid       bool     `json:"valid"`
	FailureCode string   `json:"failure_code,omitempty"`
}

func Check(req Request) Result {
	symbol := strings.ToUpper(strings.TrimSpace(req.Symbol))
	side := strings.ToLower(strings.TrimSpace(req.Side))
	out := Result{
		CheckID:    uuid.NewString(),
		Violations: make([]string, 0, 4),
	}

	if symbol == "" || req.Qty <= 0 || req.Price <= 0 {
		out.Valid = false
		out.FailureCode = "invalid_payload"
		return out
	}
	if side != "buy" && side != "sell" {
		out.Valid = false
		out.FailureCode = "invalid_side"
		return out
	}

	notional := req.Qty * req.Price
	out.Notional = notional

	normalizedRestricted := make([]string, 0, len(req.RestrictedSymbols))
	for _, s := range req.RestrictedSymbols {
		normalizedRestricted = append(normalizedRestricted, strings.ToUpper(strings.TrimSpace(s)))
	}
	if slices.Contains(normalizedRestricted, symbol) {
		out.Violations = append(out.Violations, "restricted_symbol")
	}
	if req.MaxNotional > 0 && notional > req.MaxNotional {
		out.Violations = append(out.Violations, "max_notional_exceeded")
	}

	projected := req.CurrentPosition
	if side == "buy" {
		projected += req.Qty
	} else {
		projected -= req.Qty
	}
	if req.MaxPosition > 0 && math.Abs(projected) > req.MaxPosition {
		out.Violations = append(out.Violations, "position_limit_exceeded")
	}

	out.Valid = true
	out.Allowed = len(out.Violations) == 0
	return out
}

