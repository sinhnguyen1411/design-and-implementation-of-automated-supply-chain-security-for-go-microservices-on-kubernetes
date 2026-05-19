package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/search-service/internal/health"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/search-service/internal/search"
)

var catalog = []search.Symbol{
	{Ticker: "AAPL", Name: "Apple Inc", Exchange: "NASDAQ", Sector: "Technology"},
	{Ticker: "GOOGL", Name: "Alphabet Inc", Exchange: "NASDAQ", Sector: "Technology"},
	{Ticker: "JPM", Name: "JPMorgan Chase", Exchange: "NYSE", Sector: "Finance"},
}

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(health.Status()))
	})
	http.HandleFunc("/search/symbols", func(w http.ResponseWriter, r *http.Request) {
		results := search.Query(catalog, r.URL.Query().Get("q"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(results)
	})
	fmt.Printf("search-service listening on :%s\n", port)
	_ = http.ListenAndServe(":"+port, nil)
}
