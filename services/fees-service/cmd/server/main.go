package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/fees-service/internal/fees"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/fees-service/internal/health"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(health.Status()))
	})
	http.HandleFunc("/fees/calculate", func(w http.ResponseWriter, r *http.Request) {
		result := fees.Calculate(10000, fees.TierStandard)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result)
	})
	fmt.Printf("fees-service listening on :%s\n", port)
	_ = http.ListenAndServe(":"+port, nil)
}
