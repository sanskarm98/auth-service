# Build Stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code
COPY . .

# Build the binary from the correct path
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/api

# Final Stage
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/auth-service .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./auth-service"]
