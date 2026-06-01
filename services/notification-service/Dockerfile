# syntax=docker/dockerfile:1.7
FROM golang:1.25.10-bookworm AS builder
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/service ./cmd/server

# Distroless nonroot runtime: ships CA certificates, has no shell or package
# manager (smaller CVE surface), and runs as uid 65532 by default.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/service /service
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/service"]
