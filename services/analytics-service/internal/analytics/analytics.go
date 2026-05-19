package analytics

import "math"

type DailyReturn struct {
	Date   string  `json:"date"`
	Return float64 `json:"return"`
}

type Metrics struct {
	Sharpe      float64 `json:"sharpe"`
	Sortino     float64 `json:"sortino"`
	CAGR        float64 `json:"cagr"`
	MaxDrawdown float64 `json:"max_drawdown"`
	TotalReturn float64 `json:"total_return"`
	Volatility  float64 `json:"volatility"`
}

const (
	tradingDaysPerYear = 252
	riskFreeDaily      = 0.04 / 252
)

func mean(rs []float64) float64 {
	if len(rs) == 0 {
		return 0
	}
	s := 0.0
	for _, r := range rs {
		s += r
	}
	return s / float64(len(rs))
}

func stddev(rs []float64, m float64) float64 {
	if len(rs) < 2 {
		return 0
	}
	variance := 0.0
	for _, r := range rs {
		d := r - m
		variance += d * d
	}
	return math.Sqrt(variance / float64(len(rs)-1))
}

func downsideStddev(rs []float64, target float64) float64 {
	if len(rs) < 2 {
		return 0
	}
	sum := 0.0
	n := 0
	for _, r := range rs {
		if r < target {
			d := r - target
			sum += d * d
			n++
		}
	}
	if n < 2 {
		return 0
	}
	return math.Sqrt(sum / float64(n-1))
}

func ComputeMetrics(returns []DailyReturn) Metrics {
	if len(returns) == 0 {
		return Metrics{}
	}
	raw := make([]float64, len(returns))
	for i, r := range returns {
		raw[i] = r.Return
	}
	m := mean(raw)
	vol := stddev(raw, m)
	downVol := downsideStddev(raw, riskFreeDaily)

	sharpe := 0.0
	if vol > 0 {
		sharpe = (m - riskFreeDaily) / vol * math.Sqrt(tradingDaysPerYear)
	}
	sortino := 0.0
	if downVol > 0 {
		sortino = (m - riskFreeDaily) / downVol * math.Sqrt(tradingDaysPerYear)
	}

	wealth := 1.0
	for _, r := range raw {
		wealth *= (1 + r)
	}
	n := float64(len(raw))
	cagr := math.Pow(wealth, tradingDaysPerYear/n) - 1

	peak := 1.0
	cur := 1.0
	maxDD := 0.0
	for _, r := range raw {
		cur *= (1 + r)
		if cur > peak {
			peak = cur
		}
		if dd := (peak - cur) / peak; dd > maxDD {
			maxDD = dd
		}
	}

	return Metrics{
		Sharpe:      sharpe,
		Sortino:     sortino,
		CAGR:        cagr,
		MaxDrawdown: maxDD,
		TotalReturn: wealth - 1,
		Volatility:  vol * math.Sqrt(tradingDaysPerYear),
	}
}
