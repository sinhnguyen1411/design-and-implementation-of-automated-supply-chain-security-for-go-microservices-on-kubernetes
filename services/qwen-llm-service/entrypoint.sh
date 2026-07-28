#!/bin/bash
# Start the Go Proxy in the background
/usr/local/bin/llm-proxy &

# Start the actual Ollama engine in the foreground
exec ollama serve
