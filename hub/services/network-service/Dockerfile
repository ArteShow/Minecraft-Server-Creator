FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o network-service ./cmd/network-service/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/network-service .

EXPOSE 8011
CMD ["./network-service"]
