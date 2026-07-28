# Stage 1: Build the Go proxy
FROM golang:1.25.11-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN go build -o /llm-proxy main.go

# Stage 2: The final AI Immutable Image
FROM ollama/ollama:latest

# Copy the Go proxy (which satisfies the monorepo's structural requirements)
COPY --from=builder /llm-proxy /usr/local/bin/llm-proxy

# SECURITY CRITICAL: Bake the model weights into the Docker Image during build time.
# This ensures Kyverno and Cosign can sign the EXACT model weights, 
# preventing runtime tampering or downloading of malicious GGUF files.
RUN nohup bash -c "ollama serve &" && \
    sleep 5 && \
    ollama pull qwen2.5:3b && \
    pkill ollama || true

EXPOSE 11434 8080

# Copy custom entrypoint to start both the Go proxy and Ollama
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
