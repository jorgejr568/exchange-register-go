FROM golang:1.23-alpine3.16 AS builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o exchange.bin

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /build/exchange.bin .
COPY --from=builder /build/server/static /app/static

EXPOSE 8080
CMD ["./exchange.bin", "service", "--sync", "--port", "8080"]
