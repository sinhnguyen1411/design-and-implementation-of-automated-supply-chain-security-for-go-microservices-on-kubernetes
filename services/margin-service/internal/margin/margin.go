package margin

type Req struct {
	Symbol      string  `json:"symbol"`
	Qty         float64 `json:"qty"`
	Price       float64 `json:"price"`
	MarginRatio float64 `json:"margin_ratio"`
}

type Result struct {
	Notional          float64 `json:"notional"`
	InitialMargin     float64 `json:"initial_margin"`
	MaintenanceMargin float64 `json:"maintenance_margin"`
	ExcessLiquidity   float64 `json:"excess_liquidity"`
}

func Calculate(req Req, accountEquity float64) Result {
	notional := req.Qty * req.Price
	initial := notional * req.MarginRatio
	return Result{
		Notional:          notional,
		InitialMargin:     initial,
		MaintenanceMargin: initial * 0.75,
		ExcessLiquidity:   accountEquity - initial,
	}
}
