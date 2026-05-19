package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/gateway-service/internal/gateway"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/gateway-service/internal/health"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	reg := gateway.NewRegistry()
	reg.Register(gateway.Route{Path: "/api/users", Method: "GET", Backend: "user-service", RateLimitRPS: 100})
	reg.Register(gateway.Route{Path: "/api/orders", Method: "POST", Backend: "order-service", RateLimitRPS: 50})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(health.Status()))
	})
	http.HandleFunc("/gateway/routes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(reg.List())
	})
	fmt.Printf("gateway-service listening on :%s\n", port)
	_ = http.ListenAndServe(":"+port, nil)
}
