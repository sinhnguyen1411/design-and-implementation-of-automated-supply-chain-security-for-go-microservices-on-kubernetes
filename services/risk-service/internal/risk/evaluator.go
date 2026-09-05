package risk

type Request struct {
	ExposureValue float64 `json:"exposure_value"`
	PortfolioNAV  float64 `json:"portfolio_nav"`
	VaR95         float64 `json:"var_95"`
	MaxDrawdown   float64 `json:"max_drawdown"`
	Leverage      float64 `json:"leverage"`
}

type Result struct {
	Score        int      `json:"score"`
	Band         string   `json:"band"`
	Warnings     []string `json:"warnings"`
	TradingBlock bool     `json:"trading_block"`
}

func Evaluate(req Request) Result {
	score := 0
	warnings := make([]string, 0, 4)

	concentration := 0.0
	if req.PortfolioNAV > 0 {
		concentration = req.ExposureValue / req.PortfolioNAV
	}

	if concentration > 0.4 {
		score += 30
		warnings = append(warnings, "high_single_name_exposure")
	}
	if req.VaR95 > 0.03 {
		score += 30
		warnings = append(warnings, "var95_breach")
	}
	if req.MaxDrawdown > 0.15 {
		score += 20
		warnings = append(warnings, "drawdown_breach")
	}
	if req.Leverage > 2.0 {
		score += 20
		warnings = append(warnings, "leverage_breach")
	}

	band := "low"
	block := false
	switch {
	case score >= 70:
		band = "high"
		block = true
	case score >= 40:
		band = "medium"
	}

	return Result{
		Score:        score,
		Band:         band,
		Warnings:     warnings,
		TradingBlock: block,
	}
}
