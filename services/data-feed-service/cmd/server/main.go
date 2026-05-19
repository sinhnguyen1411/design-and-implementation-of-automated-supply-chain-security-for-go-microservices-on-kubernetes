package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/data-feed-service/internal/datafeed"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/data-feed-service/internal/health"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(health.Status()))
	})
	http.HandleFunc("/datafeed/tick", func(w http.ResponseWriter, r *http.Request) {
		ticks := []datafeed.Tick{datafeed.NewTick("AAPL", 150, 1000), datafeed.NewTick("AAPL", 152, 800)}
		bar := datafeed.AggregateToBar(ticks)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(bar)
	})
	fmt.Printf("data-feed-service listening on :%s\n", port)
	_ = http.ListenAndServe(":"+port, nil)
}
