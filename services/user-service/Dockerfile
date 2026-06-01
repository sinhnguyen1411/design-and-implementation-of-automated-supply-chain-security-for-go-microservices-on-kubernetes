# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.25.10

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-bookworm AS builder
LABEL org.opencontainers.image.source="https://github.com/sinhnguyen1411/stock-trading-be"
WORKDIR /src

ENV CGO_ENABLED=0 \
    GOENV=/src/.goenv \
    GOMODCACHE=/go/pkg/mod \
    GOCACHE=/go/build-cache

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/go/build-cache \
    go mod download

COPY . .

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

# Distroless nonroot runtime: includes CA certificates, has no shell or package
# manager (smaller CVE surface than the previous Ubuntu 18.04 LTS base), and runs
# as uid 65532 by default for least privilege.
FROM gcr.io/distroless/static-debian12:nonroot AS runtime
LABEL org.opencontainers.image.title="stock-trading-user-service" \
      org.opencontainers.image.vendor="stock-trading" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.source="https://github.com/sinhnguyen1411/stock-trading-be"

COPY --from=builder /out/user-service /usr/local/bin/user-service
COPY cmd/server/config/local.yaml /etc/stock-trading/config.yaml

ENV STOCK_TRADING_CONFIG=/etc/stock-trading/config.yaml

EXPOSE 18080 19090
USER 65532:65532

ENTRYPOINT ["/usr/local/bin/user-service"]
CMD ["server", "--config", "/etc/stock-trading/config.yaml"]
