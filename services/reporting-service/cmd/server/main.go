package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/reporting-service/internal/health"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/reporting-service/internal/report"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(health.Status()))
	})
	http.HandleFunc("/reports/statement", func(w http.ResponseWriter, r *http.Request) {
		stmt := report.BuildStatement("acc-demo", []report.LineItem{{Date: "2026-01-01", Description: "initial deposit", Amount: 10000}}, 0)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(stmt)
	})
	fmt.Printf("reporting-service listening on :%s\n", port)
	_ = http.ListenAndServe(":"+port, nil)
}
