# =============================================================================
# Build stage
# =============================================================================
FROM golang:1.25.0-alpine3.20 AS builder

# Install build dependencies and security tools
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    upx \
    && update-ca-certificates

# Create non-root user for build process
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /build

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum ./

# Download and verify dependencies (cached layer)
RUN go mod download && \
    go mod verify

# Copy source code
COPY . .

# Set build arguments for versioning and metadata
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Build the application with optimizations
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT} -extldflags '-static'" \
    -a \
    -installsuffix cgo \
    -tags netgo \
    -o app \
    ./cmd/server/main.go

# Compress binary with UPX (optional, uncomment if needed)
# RUN upx --best --lzma app

# Verify the binary
RUN ./app --version || echo "Binary built successfully"

# =============================================================================
# Final stage - Distroless for maximum security
# =============================================================================
FROM gcr.io/distroless/static-debian12:nonroot AS production

# Copy timezone data and certificates from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /build/app /app

# Use distroless nonroot user (uid=65532, gid=65532)
USER nonroot:nonroot

# Set environment variables
ENV TZ=UTC
ENV GIN_MODE=release
ENV PORT=8080
ENV GOMAXPROCS=0
ENV GOMEMLIMIT=0

# Expose port
EXPOSE 8080

# Add labels for better container management
LABEL maintainer="LMS Team" \
      version="${VERSION}" \
      description="LMS Application" \
      org.opencontainers.image.source="https://github.com/your-org/lms" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.created="${BUILD_TIME}" \
      org.opencontainers.image.base.name="gcr.io/distroless/static-debian12:nonroot"

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/app", "--health-check"]

# Run the application
ENTRYPOINT ["/app"]

# =============================================================================
# Development stage (optional)
# =============================================================================
FROM golang:1.25.0-alpine3.20 AS development

# Install development tools
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    curl \
    bash \
    make

# Install air for hot reloading
RUN go install github.com/cosmtrek/air@latest

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy go files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Change ownership
RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

# Use air for development
CMD ["air"]