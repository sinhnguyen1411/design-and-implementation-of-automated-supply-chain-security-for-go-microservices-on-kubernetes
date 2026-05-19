package backtest

type Result struct {
	TotalReturn  float64 `json:"total_return"`
	MaxDrawdown  float64 `json:"max_drawdown"`
	FinalCapital float64 `json:"final_capital"`
	Trades       int     `json:"trades"`
}

func Run(prices []float64, capital float64) Result {
	if len(prices) < 2 {
		return Result{FinalCapital: capital}
	}
	peak, maxDD, trades, cur := capital, 0.0, 0, capital
	for i := 1; i < len(prices); i++ {
		ret := (prices[i] - prices[i-1]) / prices[i-1]
		if ret > 0 { cur *= (1 + ret); trades++ }
		if cur > peak { peak = cur }
		if dd := (peak - cur) / peak; dd > maxDD { maxDD = dd }
	}
	return Result{TotalReturn: (cur - capital) / capital, MaxDrawdown: maxDD, FinalCapital: cur, Trades: trades}
}
