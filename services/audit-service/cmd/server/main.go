package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/audit-service/internal/audit"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/audit-service/internal/health"
)

var auditLog = audit.NewEventLog()

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(health.Status()))
	})
	http.HandleFunc("/events/append", appendHandler)
	http.HandleFunc("/events/replay", replayHandler)
	http.HandleFunc("/events/snapshot", snapshotHandler)
	http.ListenAndServe(":"+port, nil)
}

func appendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		AggregateID string          `json:"aggregate_id"`
		Actor       string          `json:"actor"`
		Type        audit.EventType `json:"type"`
		Payload     json.RawMessage `json:"payload"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	e, err := auditLog.Append(req.AggregateID, req.Actor, req.Type, req.Payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

func replayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	aggID := r.URL.Query().Get("aggregate_id")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(auditLog.Replay(aggID))
}

func snapshotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		AggregateID string          `json:"aggregate_id"`
		State       json.RawMessage `json:"state"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	snap := auditLog.TakeSnapshot(req.AggregateID, req.State)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snap)
}
