# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o cashout \
    ./cmd/server/main.go

# Build the web server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o cashout-web \
    ./cmd/web/main.go

# Build the migration tool
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o migrate \
    ./cmd/migrate/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 -S cashout && \
    adduser -u 1000 -S cashout -G cashout

# Set working directory
WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/cashout /app/cashout
COPY --from=builder /app/cashout-web /app/cashout-web
COPY --from=builder /app/migrate /app/migrate

# Change ownership
RUN chown -R cashout:cashout /app

# Switch to non-root user
USER cashout

# Expose ports
EXPOSE 8080 8081

# Set entrypoint (default to bot, can be overridden)
ENTRYPOINT ["/app/cashout"]
