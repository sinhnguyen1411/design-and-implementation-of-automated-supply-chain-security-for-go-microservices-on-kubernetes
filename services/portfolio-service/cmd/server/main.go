package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/portfolio-service/internal/health"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(health.Status()))
	})
	fmt.Printf("portfolio-service listening on :%s\n", port)
	_ = http.ListenAndServe(":"+port, nil)
}