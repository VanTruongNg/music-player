# syntax=docker/dockerfile:1.6
# =============================================================================
# BigTech Production Dockerfile - API Gateway
# Cloud-Native, Security-First, High-Performance
# =============================================================================

# -----------------------------------------------------------------------------
# Build Stage: Advanced optimization with multi-arch support
# -----------------------------------------------------------------------------
FROM golang:1.24.2-alpine3.20 AS builder

# Build arguments for flexibility
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=dev
ARG COMMIT_SHA=unknown
ARG BUILD_TIME

# Security: Create consistent non-root user
RUN addgroup -g 65532 -S nonroot && \
    adduser -u 65532 -S nonroot -G nonroot

# Performance: Optimized dependency installation
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    upx \
    && rm -rf /var/cache/apk/* \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build
RUN chown nonroot:nonroot /build

# Performance: Dependency management with optimal caching
COPY --chown=nonroot:nonroot go.mod go.sum ./

# Security: Switch to non-root user early
USER nonroot:nonroot

RUN go mod download && go mod verify

# Source code
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

# Performance: Binary compression
RUN upx --best --lzma main || echo "UPX compression failed, continuing..."

# -----------------------------------------------------------------------------
# Production Stage: Ultra-secure, minimal runtime
# -----------------------------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot

# Build arguments for metadata
ARG VERSION=dev
ARG BUILD_TIME

# OCI metadata with comprehensive information
LABEL \
    org.opencontainers.image.title="API Gateway" \
    org.opencontainers.image.description="BigTech-grade API Gateway Microservice" \
    org.opencontainers.image.vendor="SupaGoodSongs" \
    org.opencontainers.image.version="${VERSION}" \
    org.opencontainers.image.created="${BUILD_TIME}" \
    org.opencontainers.image.source="https://github.com/VanTruongNg/music-player" \
    org.opencontainers.image.documentation="https://docs.music-player.com/gateway"

# Security: Essential runtime files
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Application files with correct ownership
COPY --from=builder --chown=65532:65532 /build/main /app/main
COPY --from=builder --chown=65532:65532 /build/configs /app/configs

WORKDIR /app

# Runtime optimization
ENV TZ=UTC \
    GIN_MODE=release \
    GOMAXPROCS=1 \
    GOMEMLIMIT=200MiB

# Security: Non-root execution
USER 65532:65532

ENTRYPOINT ["/app/main"]