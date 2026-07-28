package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// A simple Go proxy to satisfy the Go CI pipeline requirements
// and potentially add custom routing/metrics in the future before hitting Ollama.
func main() {
	fmt.Println("Starting Qwen LLM Proxy on :8080")
	fmt.Println("Forwarding requests to Ollama on :11434")

	ollamaURL, _ := url.Parse("http://127.0.0.1:11434")
	proxy := httputil.NewSingleHostReverseProxy(ollamaURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
