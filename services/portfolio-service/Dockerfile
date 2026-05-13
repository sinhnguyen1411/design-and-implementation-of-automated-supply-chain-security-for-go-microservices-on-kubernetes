FROM golang:1.25.10 AS builder
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/service ./cmd/server

FROM gcr.io/distroless/static-debian12
COPY --from=builder /out/service /service
EXPOSE 8080
ENTRYPOINT ["/service"]

