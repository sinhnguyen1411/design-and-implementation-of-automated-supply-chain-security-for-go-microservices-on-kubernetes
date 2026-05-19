package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/account-service/internal/brokerage"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/account-service/internal/health"
)

var (
	mu       sync.Mutex
	accounts = make(map[string]*brokerage.BrokerageAccount)
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }

	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(health.Status()))
	})

	http.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
		var req struct { OwnerID string `json:"owner_id"`; Type brokerage.AccountType `json:"type"` }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "bad request", http.StatusBadRequest); return }
		acc := brokerage.New(req.OwnerID, req.Type)
		mu.Lock(); accounts[acc.ID] = &acc; mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(acc)
	})

	http.HandleFunc("/accounts/deposit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
		var req struct { AccountID string `json:"account_id"`; Amount float64 `json:"amount"` }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "bad request", http.StatusBadRequest); return }
		mu.Lock(); acc, ok := accounts[req.AccountID]; mu.Unlock()
		if !ok { http.Error(w, "account not found", http.StatusNotFound); return }
		tx, err := acc.Deposit(req.Amount)
		if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(tx)
	})

	http.HandleFunc("/accounts/withdraw", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
		var req struct { AccountID string `json:"account_id"`; Amount float64 `json:"amount"` }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "bad request", http.StatusBadRequest); return }
		mu.Lock(); acc, ok := accounts[req.AccountID]; mu.Unlock()
		if !ok { http.Error(w, "account not found", http.StatusNotFound); return }
		tx, err := acc.Withdraw(req.Amount)
		if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(tx)
	})

	fmt.Printf("account-service listening on :%s\n", port)
	_ = http.ListenAndServe(":"+port, nil)
}
