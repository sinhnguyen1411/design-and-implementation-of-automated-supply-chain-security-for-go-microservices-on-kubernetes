package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/account-service/internal/account"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/account-service/internal/health"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(health.Status()))
	})
	http.HandleFunc("/account/status", func(w http.ResponseWriter, r *http.Request) {
		ownerID := r.URL.Query().Get("owner_id")
		if ownerID == "" { ownerID = "anonymous" }
		a := account.New(ownerID)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(a)
	})
	fmt.Printf("account-service listening on :%s\n", port)
	_ = http.ListenAndServe(":"+port, nil)
}
