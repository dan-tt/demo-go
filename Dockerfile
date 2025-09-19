# Multi-stage Dockerfile for Go API - Production Ready
# Stage 1: Build stage
FROM golang:1.21-alpine AS builder

# Build arguments for CI/CD pipeline
ARG GO_VERSION=1.21
ARG BUILD_TIME
ARG GIT_COMMIT
ARG GIT_BRANCH

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies with verification
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with optimization and build info
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME} -X main.Branch=${GIT_BRANCH}" \
    -a -installsuffix cgo \
    -o main cmd/server/main.go

# Stage 2: Production stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check for container orchestration
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set production environment variables
ENV GIN_MODE=release
ENV PORT=8080

# Run the application
CMD ["./main"]

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./server"]
