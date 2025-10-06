# ---- Build Stage ----
FROM golang:1.25-alpine AS builder

# Set environment variables for Go
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go binary
RUN go build -o server ./cmd/server/main.go

# ---- Final Stage ----
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the compiled binary from builder stage
COPY --from=builder /app/server .

# Copy config files
COPY --from=builder /app/config ./config

# Create storage directory
RUN mkdir -p /app/storage

# Expose port your app listens on
EXPOSE 8080

# Command to run the executable
CMD ["./server"]
