package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/kyc-service/internal/health"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/kyc-service/internal/kyc"
)

var registry = kyc.NewRegistry()

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(health.Status()))
	})
	http.HandleFunc("/kyc/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			UserID       string `json:"user_id"`
			Jurisdiction string `json:"jurisdiction"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		p := registry.Create(req.UserID, req.Jurisdiction)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	})
	http.HandleFunc("/kyc/document", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			UserID string      `json:"user_id"`
			Type   kyc.DocType `json:"type"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		doc, err := registry.AddDocument(req.UserID, req.Type, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(doc)
	})
	http.HandleFunc("/kyc/approve", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			UserID string `json:"user_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := registry.Approve(req.UserID); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	http.HandleFunc("/kyc/reject", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			UserID string `json:"user_id"`
			Reason string `json:"reason"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := registry.Reject(req.UserID, req.Reason); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	http.HandleFunc("/kyc/aml/flag", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			UserID   string        `json:"user_id"`
			Code     string        `json:"code"`
			Notes    string        `json:"notes"`
			Severity kyc.RiskLevel `json:"severity"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := registry.FlagAML(req.UserID, req.Code, req.Notes, req.Severity); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
	http.HandleFunc("/kyc/profile", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		p, err := registry.Get(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	})
	http.ListenAndServe(":"+port, nil)
}
