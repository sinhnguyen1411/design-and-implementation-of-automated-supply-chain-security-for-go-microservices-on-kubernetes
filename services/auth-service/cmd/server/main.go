package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/auth-service/internal/health"
	"github.com/sinhnguyen1411/design-and-implementation-of-automated-supply-chain-security-for-go-microservices-on-kubernetes/services/auth-service/internal/oauth2"
)

var (
	mu         sync.Mutex
	tokenStore = make(map[string]oauth2.AccessToken)
	codeStore  = make(map[string]*oauth2.AuthCode)
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }

	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(health.Status()))
	})

	http.HandleFunc("/oauth2/authorize", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		code := oauth2.IssueCode(q.Get("client_id"), q.Get("redirect_uri"), q.Get("subject"), []string{"read", "trade"})
		mu.Lock(); codeStore[code.Code] = &code; mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(code)
	})

	http.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		rawCode := r.FormValue("code")
		clientID := r.FormValue("client_id")
		mu.Lock(); code, ok := codeStore[rawCode]; mu.Unlock()
		if !ok { http.Error(w, "invalid code", http.StatusBadRequest); return }
		tok, err := oauth2.ExchangeCode(code, clientID)
		if err != nil { http.Error(w, err.Error(), http.StatusUnauthorized); return }
		mu.Lock(); tokenStore[tok.Token] = tok; mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(tok)
	})

	http.HandleFunc("/oauth2/introspect", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		clientID := r.URL.Query().Get("client_id")
		mu.Lock(); result := oauth2.Introspect(token, clientID, tokenStore); mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result)
	})

	fmt.Printf("auth-service listening on :%s\n", port)
	_ = http.ListenAndServe(":"+port, nil)
}
