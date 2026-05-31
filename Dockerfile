# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application statically
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o repo-lyzer .

# Runtime stage
FROM alpine:3.19

# Add CA certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user and group
RUN addgroup -S repolyzer && adduser -S repolyzer -G repolyzer

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/repo-lyzer .

# Create data directory and set permissions
RUN mkdir -p /app/data && chown -R repolyzer:repolyzer /app

# Use non-root user
USER repolyzer

# Set environment variables for config
ENV REPO_LYZER_CONFIG_PATH=/app/data/settings.json

# Add healthcheck to ensure daemon is running
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD pgrep repo-lyzer || exit 1

# Default command runs the daemon
CMD ["./repo-lyzer", "daemon"]
