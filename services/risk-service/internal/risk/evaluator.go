package risk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

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
	AIAnalysis   string   `json:"ai_analysis,omitempty"`
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
		AIAnalysis:   callOllamaForRisk(req),
	}
}

// callOllamaForRisk calls the qwen-llm-service to get an AI assessment
func callOllamaForRisk(req Request) string {
	prompt := fmt.Sprintf("Phân tích hồ sơ rủi ro giao dịch này để tìm điểm bất thường: Exposure: %.2f, NAV: %.2f, VaR95: %.2f, Drawdown: %.2f, Leverage: %.2f. Trả lời ngắn gọn dưới 3 câu bằng tiếng Việt.", req.ExposureValue, req.PortfolioNAV, req.VaR95, req.MaxDrawdown, req.Leverage)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model":  "qwen2.5:3b",
		"prompt": prompt,
		"stream": false,
	})

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Post("http://qwen-llm-service:11434/api/generate", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "AI Service unavailable"
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "AI Service error"
	}

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "Failed to parse AI response"
	}

	if response, ok := res["response"].(string); ok {
		return response
	}
	return "No response from AI"
}
