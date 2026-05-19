package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/alert-service/internal/alert"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/alert-service/internal/health"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(health.Status()))
	})
	http.HandleFunc("/alerts/trigger", func(w http.ResponseWriter, r *http.Request) {
		result := alert.Evaluate(alert.Alert{SymbolID: "AAPL", Threshold: 150, Type: alert.AlertAbove}, 155.5)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result)
	})
	fmt.Printf("alert-service listening on :%s\n", port)
	_ = http.ListenAndServe(":"+port, nil)
}
