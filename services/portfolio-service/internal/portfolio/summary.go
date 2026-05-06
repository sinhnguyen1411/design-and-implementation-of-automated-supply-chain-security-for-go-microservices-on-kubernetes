package portfolio

type Position struct {
	Symbol string  `json:"symbol"`
	Qty    float64 `json:"qty"`
	Price  float64 `json:"price"`
}

type Summary struct {
	Cash                float64            `json:"cash"`
	MarketValue         float64            `json:"market_value"`
	TotalEquity         float64            `json:"total_equity"`
	ConcentrationRatio  float64            `json:"concentration_ratio"`
	SectorDiversifyHint string             `json:"sector_diversify_hint"`
	Weights             map[string]float64 `json:"weights"`
}

func BuildSummary(cash float64, positions []Position) Summary {
	weights := make(map[string]float64, len(positions))
	marketValue := 0.0
	largestValue := 0.0

	for _, p := range positions {
		value := p.Qty * p.Price
		marketValue += value
		if value > largestValue {
			largestValue = value
		}
		weights[p.Symbol] = value
	}

	total := cash + marketValue
	if marketValue > 0 {
		for symbol, value := range weights {
			weights[symbol] = value / marketValue
		}
	}

	concentration := 0.0
	if marketValue > 0 {
		concentration = largestValue / marketValue
	}

	hint := "balanced"
	if concentration >= 0.5 {
		hint = "high_concentration"
	} else if concentration >= 0.3 {
		hint = "moderate_concentration"
	}

	return Summary{
		Cash:                cash,
		MarketValue:         marketValue,
		TotalEquity:         total,
		ConcentrationRatio:  concentration,
		SectorDiversifyHint: hint,
		Weights:             weights,
	}
}
