package analytics

type Trade struct {
	Symbol    string  `json:"symbol"`
	Qty       float64 `json:"qty"`
	BuyPrice  float64 `json:"buy_price"`
	SellPrice float64 `json:"sell_price"`
}

type PnLResult struct {
	GrossPnL float64 `json:"gross_pnl"`
	NetPnL   float64 `json:"net_pnl"`
	WinRate  float64 `json:"win_rate"`
	Trades   int     `json:"trades"`
}

func ComputePnL(trades []Trade, feePct float64) PnLResult {
	if len(trades) == 0 {
		return PnLResult{}
	}
	gross, wins := 0.0, 0
	for _, t := range trades {
		pnl := (t.SellPrice - t.BuyPrice) * t.Qty
		gross += pnl
		if pnl > 0 { wins++ }
	}
	fees := gross * feePct
	return PnLResult{GrossPnL: gross, NetPnL: gross - fees, WinRate: float64(wins) / float64(len(trades)), Trades: len(trades)}
}
