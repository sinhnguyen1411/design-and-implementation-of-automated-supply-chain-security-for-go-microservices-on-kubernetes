package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/analytics-service/internal/analytics"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/analytics-service/internal/health"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(health.Status()))
	})
	http.HandleFunc("/analytics/pnl", func(w http.ResponseWriter, r *http.Request) {
		trades := []analytics.Trade{{Symbol: "AAPL", Qty: 10, BuyPrice: 150, SellPrice: 160}}
		result := analytics.ComputePnL(trades, 0.001)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result)
	})
	fmt.Printf("analytics-service listening on :%s\n", port)
	_ = http.ListenAndServe(":"+port, nil)
}
