# syntax=docker/dockerfile:1.6
# =============================================================================
# BigTech Production Dockerfile - Auth Service
# Security, Performance, and Cloud-Native Optimized
# =============================================================================

# -----------------------------------------------------------------------------
# Build Stage: Multi-arch support and advanced optimization
# -----------------------------------------------------------------------------
FROM golang:1.24.2-alpine3.20 AS builder

# Build arguments for flexibility
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=dev
ARG COMMIT_SHA=unknown
ARG BUILD_TIME

# Security: Create non-root user with specific UID/GID
RUN addgroup -g 65532 -S nonroot && \
    adduser -u 65532 -S nonroot -G nonroot

# Performance: Install build dependencies in single optimized layer
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    upx \
    && rm -rf /var/cache/apk/* \
    && rm -rf /var/lib/apt/lists/*

# Security: Set secure working directory with proper permissions
WORKDIR /build
RUN chown nonroot:nonroot /build

# Performance: Copy dependency files first for optimal layer caching
COPY --chown=nonroot:nonroot go.mod go.sum ./

# Security: Switch to non-root user for build
USER nonroot:nonroot

# Performance: Download and verify dependencies as separate cached layer
RUN go mod download && go mod verify

# Build optimization: Copy source code
COPY --chown=nonroot:nonroot . .

# Performance: Generate wire dependencies efficiently
RUN cd cmd && go run github.com/google/wire/cmd/wire

# Build optimization: Create highly optimized multi-arch binary
RUN CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    go build \
    -ldflags="-w -s -X main.version=${VERSION} -X main.commit=${COMMIT_SHA} -X main.buildTime=${BUILD_TIME} -extldflags '-static'" \
    -a \
    -installsuffix cgo \
    -tags "netgo,osusergo,static_build" \
    -trimpath \
    -o main \
    ./cmd

# Performance: Compress binary for smaller final image
RUN upx --best --lzma main || echo "UPX compression failed, continuing..."

# -----------------------------------------------------------------------------
# Production Stage: Ultra-minimal, secure runtime environment
# -----------------------------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot

# Build arguments passed to final stage
ARG VERSION=dev
ARG BUILD_TIME

# Metadata following OCI standards with build info
LABEL \
    org.opencontainers.image.title="Auth Service" \
    org.opencontainers.image.description="BigTech-grade Authentication Microservice" \
    org.opencontainers.image.vendor="SupaGoodSongs" \
    org.opencontainers.image.version="${VERSION}" \
    org.opencontainers.image.created="${BUILD_TIME}" \
    org.opencontainers.image.source="https://github.com/VanTruongNg/music-player" \
    org.opencontainers.image.licenses="MIT" \
    org.opencontainers.image.documentation="https://docs.music-player.com/auth-service"

# Security: Copy essential certificates and timezone data
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Performance: Copy optimized binary and configuration
COPY --from=builder --chown=65532:65532 /build/main /app/main
COPY --from=builder --chown=65532:65532 /build/configs /app/configs

# Security: Set working directory with proper ownership
WORKDIR /app

# Environment: Production-ready runtime configuration
ENV TZ=UTC \
    GIN_MODE=release \
    GOMAXPROCS=2 \
    GOMEMLIMIT=450MiB

# Security: Use distroless nonroot user (UID 65532)
USER 65532:65532

# Performance: Optimized container startup
ENTRYPOINT ["/app/main"]