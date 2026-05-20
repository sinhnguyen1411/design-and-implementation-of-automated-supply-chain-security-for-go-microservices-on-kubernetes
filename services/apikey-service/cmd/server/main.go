package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/apikey-service/internal/apikey"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/apikey-service/internal/health"
)

var store = apikey.NewStore()

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(health.Status()))
	})
	http.HandleFunc("/apikeys/issue", issueHandler)
	http.HandleFunc("/apikeys/validate", validateHandler)
	http.HandleFunc("/apikeys/revoke", revokeHandler)
	http.HandleFunc("/apikeys/rotate", rotateHandler)
	http.HandleFunc("/apikeys/list", listHandler)
	http.ListenAndServe(":"+port, nil)
}

func issueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		OwnerID   string         `json:"owner_id"`
		Name      string         `json:"name"`
		Scopes    []apikey.Scope `json:"scopes"`
		Tier      apikey.Tier    `json:"tier"`
		RateLimit int            `json:"rate_limit"`
		TTLHours  *int           `json:"ttl_hours,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var ttl *time.Duration
	if req.TTLHours != nil {
		d := time.Duration(*req.TTLHours) * time.Hour
		ttl = &d
	}
	ak, err := store.Issue(req.OwnerID, req.Name, req.Scopes, req.Tier, req.RateLimit, ttl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ak)
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Key   string       `json:"key"`
		Scope apikey.Scope `json:"scope"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ak, err := store.Validate(req.Key, req.Scope)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ak)
}

func revokeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		KeyID string `json:"key_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := store.Revoke(req.KeyID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func rotateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		KeyID string `json:"key_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newAK, err := store.Rotate(req.KeyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newAK)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	ownerID := r.URL.Query().Get("owner_id")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store.ListByOwner(ownerID))
}
