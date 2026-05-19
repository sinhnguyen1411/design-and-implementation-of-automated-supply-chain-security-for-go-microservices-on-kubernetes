package alert

import "fmt"

type AlertType string

const (
	AlertAbove AlertType = "above"
	AlertBelow AlertType = "below"
)

type Alert struct {
	SymbolID  string    `json:"symbol_id"`
	Threshold float64   `json:"threshold"`
	Type      AlertType `json:"type"`
}

type TriggerResult struct {
	Fired   bool   `json:"fired"`
	Message string `json:"message"`
}

func Evaluate(a Alert, currentPrice float64) TriggerResult {
	switch a.Type {
	case AlertAbove:
		if currentPrice > a.Threshold {
			return TriggerResult{Fired: true, Message: fmt.Sprintf("%s crossed above %.2f (current: %.2f)", a.SymbolID, a.Threshold, currentPrice)}
		}
	case AlertBelow:
		if currentPrice < a.Threshold {
			return TriggerResult{Fired: true, Message: fmt.Sprintf("%s dropped below %.2f (current: %.2f)", a.SymbolID, a.Threshold, currentPrice)}
		}
	}
	return TriggerResult{Fired: false, Message: "no trigger"}
}
