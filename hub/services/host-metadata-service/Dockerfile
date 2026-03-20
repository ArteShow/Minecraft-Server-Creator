FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o host-metadata-service ./cmd/host-metadata-service/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/host-metadata-service .

EXPOSE 8012
CMD ["./host-metadata-service"]
