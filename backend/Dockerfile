# Build Stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy dependency files first for Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o kite-server ./cmd/server/main.go

# Production Stage
FROM alpine:latest
WORKDIR /app

# We need the migrations directory at runtime
COPY --from=builder /app/db/migrations ./db/migrations
COPY --from=builder /app/kite-server .

# Expose the API port
EXPOSE 8080

CMD ["./kite-server"]