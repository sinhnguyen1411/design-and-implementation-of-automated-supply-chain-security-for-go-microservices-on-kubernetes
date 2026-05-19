package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/backtest-service/internal/backtest"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/backtest-service/internal/health"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(health.Status()))
	})
	http.HandleFunc("/backtest/run", func(w http.ResponseWriter, r *http.Request) {
		result := backtest.Run([]float64{100, 102, 99, 105, 110, 108, 115}, 100000)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result)
	})
	fmt.Printf("backtest-service listening on :%s\n", port)
	_ = http.ListenAndServe(":"+port, nil)
}
