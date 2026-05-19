package backtest

import "math"

type Signal int

const (
	Hold Signal = iota
	Buy
	Sell
)

type Strategy interface {
	Name() string
	Signal(prices []float64, idx int) Signal
}

type SMAStrategy struct {
	Short int
	Long  int
}

func (s *SMAStrategy) Name() string { return "sma" }

func sma(prices []float64, end, period int) float64 {
	if end < period {
		return 0
	}
	sum := 0.0
	for i := end - period; i < end; i++ {
		sum += prices[i]
	}
	return sum / float64(period)
}

func (s *SMAStrategy) Signal(prices []float64, idx int) Signal {
	if idx < s.Long {
		return Hold
	}
	prevShort := sma(prices, idx, s.Short)
	prevLong := sma(prices, idx, s.Long)
	curShort := sma(prices, idx+1, s.Short)
	curLong := sma(prices, idx+1, s.Long)
	if prevShort <= prevLong && curShort > curLong {
		return Buy
	}
	if prevShort >= prevLong && curShort < curLong {
		return Sell
	}
	return Hold
}

type MomentumStrategy struct {
	Window    int
	Threshold float64
}

func (m *MomentumStrategy) Name() string { return "momentum" }

func (m *MomentumStrategy) Signal(prices []float64, idx int) Signal {
	if idx < m.Window {
		return Hold
	}
	ret := (prices[idx] - prices[idx-m.Window]) / prices[idx-m.Window]
	if ret > m.Threshold {
		return Buy
	}
	if ret < -m.Threshold {
		return Sell
	}
	return Hold
}

type TradeRecord struct {
	Day    int     `json:"day"`
	Action Signal  `json:"action"`
	Price  float64 `json:"price"`
}

type BacktestResult struct {
	StrategyName string        `json:"strategy_name"`
	TotalReturn  float64       `json:"total_return"`
	MaxDrawdown  float64       `json:"max_drawdown"`
	Sharpe       float64       `json:"sharpe"`
	FinalCapital float64       `json:"final_capital"`
	Trades       int           `json:"trades"`
	TradeLog     []TradeRecord `json:"trade_log"`
}

type BenchmarkResult struct {
	TotalReturn float64 `json:"total_return"`
	MaxDrawdown float64 `json:"max_drawdown"`
}

type WalkForwardResult struct {
	StrategyName string           `json:"strategy_name"`
	Folds        int              `json:"folds"`
	OutOfSample  []BacktestResult `json:"out_of_sample"`
	Benchmark    BenchmarkResult  `json:"benchmark"`
}

func RunBacktest(strategy Strategy, prices []float64, initialCapital float64) BacktestResult {
	if len(prices) < 2 {
		return BacktestResult{StrategyName: strategy.Name(), FinalCapital: initialCapital}
	}
	capital := initialCapital
	peak := capital
	maxDD := 0.0
	inPosition := false
	entryPrice := 0.0
	var trades []TradeRecord
	var dailyReturns []float64
	prev := capital

	for i := 0; i < len(prices)-1; i++ {
		sig := strategy.Signal(prices, i)
		if sig == Buy && !inPosition {
			inPosition = true
			entryPrice = prices[i+1]
			trades = append(trades, TradeRecord{Day: i + 1, Action: Buy, Price: entryPrice})
		} else if sig == Sell && inPosition {
			ret := (prices[i+1] - entryPrice) / entryPrice
			capital *= (1 + ret)
			inPosition = false
			trades = append(trades, TradeRecord{Day: i + 1, Action: Sell, Price: prices[i+1]})
		}
		if inPosition {
			curValue := capital * (prices[i+1] / prices[i])
			dailyReturns = append(dailyReturns, (curValue-prev)/prev)
			prev = curValue
			if curValue > peak {
				peak = curValue
			}
			if dd := (peak - curValue) / peak; dd > maxDD {
				maxDD = dd
			}
		}
	}

	sharpe := 0.0
	if len(dailyReturns) > 1 {
		sum, n := 0.0, float64(len(dailyReturns))
		for _, r := range dailyReturns {
			sum += r
		}
		avg := sum / n
		variance := 0.0
		for _, r := range dailyReturns {
			d := r - avg
			variance += d * d
		}
		std := math.Sqrt(variance / (n - 1))
		if std > 0 {
			sharpe = (avg - 0.04/252) / std * math.Sqrt(252)
		}
	}

	return BacktestResult{
		StrategyName: strategy.Name(),
		TotalReturn:  (capital - initialCapital) / initialCapital,
		MaxDrawdown:  maxDD,
		Sharpe:       sharpe,
		FinalCapital: capital,
		Trades:       len(trades),
		TradeLog:     trades,
	}
}

func BuyAndHold(prices []float64, initialCapital float64) BenchmarkResult {
	if len(prices) < 2 {
		return BenchmarkResult{}
	}
	totalRet := (prices[len(prices)-1] - prices[0]) / prices[0]
	capital := initialCapital
	peak := capital
	maxDD := 0.0
	for _, p := range prices[1:] {
		capital = initialCapital * (p / prices[0])
		if capital > peak {
			peak = capital
		}
		if dd := (peak - capital) / peak; dd > maxDD {
			maxDD = dd
		}
	}
	return BenchmarkResult{TotalReturn: totalRet, MaxDrawdown: maxDD}
}

func WalkForward(strategy Strategy, prices []float64, initialCapital float64, folds int) WalkForwardResult {
	if folds < 2 || len(prices) < folds*4 {
		folds = 2
	}
	foldSize := len(prices) / folds
	var oos []BacktestResult
	for i := 0; i < folds; i++ {
		start := i * foldSize
		end := start + foldSize
		if end > len(prices) {
			end = len(prices)
		}
		result := RunBacktest(strategy, prices[start:end], initialCapital)
		oos = append(oos, result)
	}
	return WalkForwardResult{
		StrategyName: strategy.Name(),
		Folds:        folds,
		OutOfSample:  oos,
		Benchmark:    BuyAndHold(prices, initialCapital),
	}
}
