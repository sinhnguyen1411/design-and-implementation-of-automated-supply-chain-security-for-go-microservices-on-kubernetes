package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/margin-service/internal/health"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/margin-service/internal/margin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(health.Status()))
	})
	http.HandleFunc("/margin/calculate", func(w http.ResponseWriter, r *http.Request) {
		result := margin.Calculate(margin.Req{Symbol: "AAPL", Qty: 100, Price: 150, MarginRatio: 0.25}, 10000)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result)
	})
	fmt.Printf("margin-service listening on :%s\n", port)
	_ = http.ListenAndServe(":"+port, nil)
}
