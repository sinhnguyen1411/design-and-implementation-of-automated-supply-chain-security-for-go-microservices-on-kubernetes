package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/watchlist-service/internal/health"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/watchlist-service/internal/watchlist"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(health.Status()))
	})
	http.HandleFunc("/watchlist/items", func(w http.ResponseWriter, r *http.Request) {
		wl := watchlist.New(r.URL.Query().Get("owner_id"))
		wl.Add("AAPL"); wl.Add("GOOGL"); wl.Add("MSFT")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(wl)
	})
	fmt.Printf("watchlist-service listening on :%s\n", port)
	_ = http.ListenAndServe(":"+port, nil)
}
