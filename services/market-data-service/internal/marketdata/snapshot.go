package marketdata

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Request struct {
	Symbol       string  `json:"symbol"`
	Bid          float64 `json:"bid"`
	Ask          float64 `json:"ask"`
	Last         float64 `json:"last"`
	TickTimeUnix int64   `json:"tick_time_unix"`
}

type Snapshot struct {
	RequestID string  `json:"request_id"`
	Symbol    string  `json:"symbol"`
	Mid       float64 `json:"mid"`
	SpreadBps float64 `json:"spread_bps"`
	Last      float64 `json:"last"`
	Stale     bool    `json:"stale"`
	Quality   string  `json:"quality"`
	Valid     bool    `json:"valid"`
	Reason    string  `json:"reason,omitempty"`
}

func BuildSnapshot(req Request) Snapshot {
	symbol := strings.ToUpper(strings.TrimSpace(req.Symbol))
	out := Snapshot{
		RequestID: uuid.NewString(),
		Symbol:    symbol,
		Last:      req.Last,
	}

	if symbol == "" || req.Bid <= 0 || req.Ask <= 0 || req.Bid > req.Ask {
		out.Valid = false
		out.Reason = "invalid_quote"
		return out
	}

	mid := (req.Bid + req.Ask) / 2
	spreadBps := ((req.Ask - req.Bid) / mid) * 10000
	stale := false
	if req.TickTimeUnix > 0 {
		stale = time.Since(time.Unix(req.TickTimeUnix, 0).UTC()) > 30*time.Second
	}

	quality := "wide"
	if spreadBps <= 5 {
		quality = "tight"
	} else if spreadBps <= 20 {
		quality = "normal"
	}

	out.Mid = mid
	out.SpreadBps = spreadBps
	out.Stale = stale
	out.Quality = quality
	out.Valid = true
	return out
}

