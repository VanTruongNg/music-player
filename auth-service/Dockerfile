# syntax=docker/dockerfile:1.6
# =============================================================================
# BigTech-grade Dockerfile (Prod) - Auth Service
# - Secure: non-root, distroless
# - Small & fast: static build, trimpath, optional UPX
# - K8s-ready: no HEALTHCHECK here (use probes later)
# =============================================================================

# -----------------------------------------------------------------------------
# Build Stage
# -----------------------------------------------------------------------------
FROM golang:1.25.0-alpine3.20 AS builder

# Build arguments
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=dev
ARG COMMIT_SHA=unknown
ARG BUILD_TIME

# Path to main package (override if you have multiple cmds, e.g. ./cmd/auth)
ARG CMD_PATH=./cmd

# Whether to compress with UPX (true/false)
ARG COMPRESS=true

# Create non-root user for build
RUN addgroup -g 65532 -S nonroot && adduser -u 65532 -S nonroot -G nonroot

# Tools for build
RUN apk add --no-cache git ca-certificates tzdata upx

WORKDIR /build
COPY --chown=nonroot:nonroot go.mod go.sum ./
USER nonroot:nonroot

# Go module warmup (cache-friendly)
RUN go mod download && go mod verify

# Copy source
COPY --chown=nonroot:nonroot . .

# Optional: generate DI with wire if used (no-op if not)
# Tip: ensure wire gen files are committed in CI to avoid network variance.
RUN [ -d "cmd" ] && cd cmd && go run github.com/google/wire/cmd/wire || true

# Build statically (CGO off) with sane flags
RUN CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    go build \
    -trimpath \
    -tags "netgo,osusergo" \
    -ldflags "-w -s -X main.version=${VERSION} -X main.commit=${COMMIT_SHA} -X main.buildTime=${BUILD_TIME}" \
    -o /build/main \
    ${CMD_PATH}

# Optional: compress to shrink image (may hinder debugging/pprof)
RUN if [ "${COMPRESS}" = "true" ]; then upx --best --lzma /build/main || true; fi

# -----------------------------------------------------------------------------
# Runtime Stage (distroless)
# -----------------------------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot

ARG VERSION=dev
ARG BUILD_TIME

# OCI metadata
LABEL org.opencontainers.image.title="Auth Service" \
    org.opencontainers.image.description="BigTech-grade Authentication Microservice" \
    org.opencontainers.image.vendor="SupaGoodSongs" \
    org.opencontainers.image.version="${VERSION}" \
    org.opencontainers.image.created="${BUILD_TIME}" \
    org.opencontainers.image.source="https://github.com/VanTruongNg/music-player" \
    org.opencontainers.image.licenses="MIT"

# Certs & timezones
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# App
WORKDIR /app
COPY --from=builder --chown=65532:65532 /build/main /app/main
# If you ship runtime configs, keep them readonly
COPY --from=builder --chown=65532:65532 /build/configs /app/configs

# Minimal runtime env; leave CPU/Mem auto-tuning to orchestrator
ENV TZ=UTC \
    GIN_MODE=release

# Distroless nonroot user (65532)
USER 65532:65532

# No EXPOSE, no HEALTHCHECK (we’ll use K8s probes later)
ENTRYPOINT ["/app/main"]
