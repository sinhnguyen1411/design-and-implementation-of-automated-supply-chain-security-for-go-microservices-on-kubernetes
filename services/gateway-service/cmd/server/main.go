package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/gateway-service/internal/gateway"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/gateway-service/internal/health"
)

var limiter = gateway.NewPerClientLimiter(100, 10)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(health.Status()))
	})
	http.HandleFunc("/ratelimit/check", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			ClientKey string `json:"client_key"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		allowed := limiter.Allow(req.ClientKey)
		status := http.StatusOK
		if !allowed {
			status = http.StatusTooManyRequests
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]bool{"allowed": allowed})
	})
	http.HandleFunc("/ratelimit/stats", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(limiter.Stats())
	})
	http.ListenAndServe(":"+port, nil)
}
