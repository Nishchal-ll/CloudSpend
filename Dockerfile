# ==============================================================================
# CloudSpend — Multi-stage Production Dockerfile
# Optimized for Azure Container Registry (ACR) & Azure App Service (Linux)
# ==============================================================================

# --- Stage 1: Build the Go application binary ---
FROM golang:alpine AS builder

WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates tzdata

# Cache Go modules layer
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Compile static binary with zero CGO dependencies
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o cloudspend main.go

# --- Stage 2: Minimal Production Image ---
FROM alpine:3.20

WORKDIR /app

# Install root SSL certificates for outgoing HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/cloudspend /app/cloudspend

# Copy HTML templates and static assets
COPY templates/ /app/templates/
COPY static/ /app/static/

# Default environment variables for Azure App Service
ENV PORT=8080
ENV GIN_MODE=release

# Expose standard container port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

# Start CloudSpend server
CMD ["/app/cloudspend"]
