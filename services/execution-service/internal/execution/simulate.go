package execution

import "github.com/google/uuid"

type Request struct {
	OrderID            string  `json:"order_id"`
	OrderQty           float64 `json:"order_qty"`
	AvailableLiquidity float64 `json:"available_liquidity"`
	LimitPrice         float64 `json:"limit_price"`
	SlippageBps        float64 `json:"slippage_bps"`
}

type Result struct {
	ExecutionID  string  `json:"execution_id"`
	Status       string  `json:"status"`
	FilledQty    float64 `json:"filled_qty"`
	RemainingQty float64 `json:"remaining_qty"`
	AveragePrice float64 `json:"average_price"`
	Valid        bool    `json:"valid"`
	Reason       string  `json:"reason,omitempty"`
}

func Simulate(req Request) Result {
	out := Result{ExecutionID: uuid.NewString(), Status: "rejected"}
	if req.OrderQty <= 0 || req.LimitPrice <= 0 {
		out.Valid = false
		out.Reason = "invalid_order"
		return out
	}

	if req.AvailableLiquidity <= 0 {
		out.Valid = true
		out.Reason = "no_liquidity"
		out.AveragePrice = req.LimitPrice
		return out
	}

	fillQty := req.OrderQty
	if req.AvailableLiquidity < req.OrderQty {
		fillQty = req.AvailableLiquidity
		out.Status = "partial_fill"
	} else {
		out.Status = "filled"
	}

	out.FilledQty = fillQty
	out.RemainingQty = req.OrderQty - fillQty
	out.AveragePrice = req.LimitPrice * (1 + req.SlippageBps/10000.0)
	out.Valid = true
	return out
}

