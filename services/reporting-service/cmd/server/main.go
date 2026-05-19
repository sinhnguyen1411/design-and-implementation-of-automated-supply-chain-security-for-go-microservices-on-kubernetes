package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/reporting-service/internal/health"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/reporting-service/internal/report"
)

var (
	mu      sync.Mutex
	ledgers = make(map[string]*report.Ledger)
)

func getLedger(accountID string) *report.Ledger {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := ledgers[accountID]; !ok {
		ledgers[accountID] = &report.Ledger{}
	}
	return ledgers[accountID]
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(health.Status()))
	})
	http.HandleFunc("/tax/acquire", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			AccountID string        `json:"account_id"`
			Lot       report.TaxLot `json:"lot"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		getLedger(req.AccountID).Acquire(req.Lot)
		w.WriteHeader(http.StatusCreated)
	})
	http.HandleFunc("/tax/sell", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			AccountID string      `json:"account_id"`
			Sale      report.Sale `json:"sale"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		gains, err := getLedger(req.AccountID).Sell(req.Sale)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(gains)
	})
	http.HandleFunc("/tax/report", func(w http.ResponseWriter, r *http.Request) {
		accountID := r.URL.Query().Get("account_id")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(getLedger(accountID).Report())
	})
	http.ListenAndServe(":"+port, nil)
}
