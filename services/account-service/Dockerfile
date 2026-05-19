ROM golang:1.25.10-bookworm AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/service ./cmd/server

FROM ubuntu:18.04
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/* \
  && groupadd -g 65532 nonroot && useradd -u 65532 -g 65532 -s /bin/false nonroot
COPY --from=builder /out/service /service
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/service"]