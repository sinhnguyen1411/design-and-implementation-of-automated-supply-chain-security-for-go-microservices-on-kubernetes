package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/alert-service/internal/health"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/alert-service/internal/pricealert"
)

var store = pricealert.NewStore()

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }

	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(health.Status()))
	})

	http.HandleFunc("/alerts/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
		var req struct {
			OwnerID   string                  `json:"owner_id"`
			Symbol    string                  `json:"symbol"`
			Threshold float64                 `json:"threshold"`
			Type      pricealert.AlertType    `json:"type"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "bad request", http.StatusBadRequest); return }
		a := store.Register(req.OwnerID, req.Symbol, req.Threshold, req.Type)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(a)
	})

	http.HandleFunc("/alerts/evaluate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
		var prices map[string]float64
		if err := json.NewDecoder(r.Body).Decode(&prices); err != nil { http.Error(w, "bad request", http.StatusBadRequest); return }
		events := store.EvaluateBatch(prices)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(events)
	})

	http.HandleFunc("/alerts/list", func(w http.ResponseWriter, r *http.Request) {
		ownerID := r.URL.Query().Get("owner_id")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(store.List(ownerID))
	})

	fmt.Printf("alert-service listening on :%s\n", port)
	_ = http.ListenAndServe(":"+port, nil)
}
