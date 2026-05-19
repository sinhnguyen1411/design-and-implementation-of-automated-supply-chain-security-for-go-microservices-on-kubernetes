package fees

type Tier string

const (
	TierStandard      Tier = "standard"
	TierPro           Tier = "pro"
	TierInstitutional Tier = "institutional"
)

type FeeResult struct {
	Notional   float64 `json:"notional"`
	Commission float64 `json:"commission"`
	TotalFee   float64 `json:"total_fee"`
	NetAmount  float64 `json:"net_amount"`
}

var tierRates = map[Tier]float64{TierStandard: 0.0010, TierPro: 0.0005, TierInstitutional: 0.0002}

func Calculate(notional float64, tier Tier) FeeResult {
	rate, ok := tierRates[tier]
	if !ok { rate = tierRates[TierStandard] }
	commission := notional * rate
	return FeeResult{Notional: notional, Commission: commission, TotalFee: commission, NetAmount: notional - commission}
}
