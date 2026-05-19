package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/backtest-service/internal/backtest"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/backtest-service/internal/health"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(health.Status()))
	})
	http.HandleFunc("/backtest/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Strategy string    `json:"strategy"`
			Prices   []float64 `json:"prices"`
			Capital  float64   `json:"capital"`
			Folds    int       `json:"folds"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.Capital <= 0 {
			req.Capital = 10000
		}
		var strategy backtest.Strategy
		switch req.Strategy {
		case "sma":
			strategy = &backtest.SMAStrategy{Short: 5, Long: 20}
		default:
			strategy = &backtest.MomentumStrategy{Window: 5, Threshold: 0.03}
		}
		var result interface{}
		if req.Folds > 1 {
			result = backtest.WalkForward(strategy, req.Prices, req.Capital, req.Folds)
		} else {
			result = backtest.RunBacktest(strategy, req.Prices, req.Capital)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
	http.HandleFunc("/backtest/benchmark", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Prices  []float64 `json:"prices"`
			Capital float64   `json:"capital"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.Capital <= 0 {
			req.Capital = 10000
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(backtest.BuyAndHold(req.Prices, req.Capital))
	})
	http.ListenAndServe(":"+port, nil)
}
