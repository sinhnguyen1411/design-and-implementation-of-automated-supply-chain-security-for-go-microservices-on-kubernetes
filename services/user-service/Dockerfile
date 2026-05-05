# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.25.9

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-bookworm AS builder
LABEL org.opencontainers.image.source="https://github.com/sinhnguyen1411/stock-trading-be"
WORKDIR /src

# Enable deterministic builds and cache modules/test data between builds
ENV CGO_ENABLED=0 \
    GOENV=/src/.goenv \
    GOMODCACHE=/go/pkg/mod \
    GOCACHE=/go/build-cache

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/go/build-cache \
    go mod download

COPY . .

# Run the unit tests inside the container build to fail early on regressions
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/go/build-cache \
    go test ./...

ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/go/build-cache \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /out/user-service .

FROM gcr.io/distroless/base-debian12:nonroot@sha256:8b9f2e503e55aff85b79d6b22c7a63a65170e8698ae80de680e3f5ea600977bf AS runtime
LABEL org.opencontainers.image.title="stock-trading-user-service" \
      org.opencontainers.image.vendor="stock-trading" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.source="https://github.com/sinhnguyen1411/stock-trading-be"

WORKDIR /app
COPY --from=builder /out/user-service /usr/local/bin/user-service
# Ship a sane default config for local runs. Kubernetes will mount its own ConfigMap/Secret.
COPY cmd/server/config/local.yaml /etc/stock-trading/config.yaml

ENV STOCK_TRADING_CONFIG=/etc/stock-trading/config.yaml

EXPOSE 18080 19090
USER nonroot:nonroot

ENTRYPOINT ["/usr/local/bin/user-service"]
CMD ["server", "--config", "/etc/stock-trading/config.yaml"]
