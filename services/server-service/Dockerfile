# ---------------------------
# Builder stage (Go binary)
# ---------------------------
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server-service ./cmd/server-service/main.go

# ---------------------------
# Runtime stage (with Java 21)
# ---------------------------
FROM eclipse-temurin:21-jdk-alpine

# Set working directory
WORKDIR /app

# Copy Go binary from builder
COPY --from=builder /app/server-service .

# Expose port
EXPOSE 8003

# Run the Go server
CMD ["./server-service"]
