FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server-service-v2 ./cmd/server-service-v2/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server-service-v2 .

EXPOSE 8003
CMD ["./server-service-v2"]
