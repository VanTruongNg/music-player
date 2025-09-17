# syntax=docker/dockerfile:1.6
# =============================================================================
# BigTech-grade Dockerfile (Prod) - API Gateway
# - Secure: non-root, distroless
# - Small & fast: static build, trimpath, optional UPX
# - K8s-ready: no HEALTHCHECK (use probes later)
# =============================================================================

# -----------------------------------------------------------------------------
# Build Stage
# -----------------------------------------------------------------------------
FROM golang:1.25.0-alpine3.20 AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=dev
ARG COMMIT_SHA=unknown
ARG BUILD_TIME
# Path to main package (override if multiple commands, e.g. ./cmd/gateway)
ARG CMD_PATH=./cmd
# Toggle UPX compression (true/false)
ARG COMPRESS=true

# Non-root for build
RUN addgroup -g 65532 -S nonroot && adduser -u 65532 -S nonroot -G nonroot

# Tooling
RUN apk add --no-cache git ca-certificates tzdata upx

WORKDIR /build
COPY --chown=nonroot:nonroot go.mod go.sum ./
USER nonroot:nonroot

# Cache deps
RUN go mod download && go mod verify

# Source
COPY --chown=nonroot:nonroot . .

# Optional: wire generate if used
RUN [ -d "cmd" ] && cd cmd && go run github.com/google/wire/cmd/wire || true

# Build static binary (CGO off), clean flags
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -trimpath \
    -tags "netgo,osusergo" \
    -ldflags "-w -s -X main.version=${VERSION} -X main.commit=${COMMIT_SHA} -X main.buildTime=${BUILD_TIME}" \
    -o /build/main \
    ${CMD_PATH}

# Optional: shrink further (can hinder debugging/pprof)
RUN if [ "${COMPRESS}" = "true" ]; then upx --best --lzma /build/main || true; fi

# -----------------------------------------------------------------------------
# Runtime Stage
# -----------------------------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot

ARG VERSION=dev
ARG BUILD_TIME

LABEL org.opencontainers.image.title="API Gateway" \
    org.opencontainers.image.description="BigTech-grade API Gateway Microservice" \
    org.opencontainers.image.vendor="SupaGoodSongs" \
    org.opencontainers.image.version="${VERSION}" \
    org.opencontainers.image.created="${BUILD_TIME}" \
    org.opencontainers.image.source="https://github.com/VanTruongNg/music-player" \
    org.opencontainers.image.licenses="MIT"

# Certs & timezones
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app
COPY --from=builder --chown=65532:65532 /build/main /app/main
# Nếu cần config runtime, mount qua volume/env thay vì bake vào image
# COPY --from=builder --chown=65532:65532 /build/configs /app/configs

# Runtime env tối thiểu; để orchestrator tự điều phối CPU/Mem
ENV TZ=UTC \
    GIN_MODE=release

# (Tuỳ chọn) Graceful shutdown signal cho orchestrator
STOPSIGNAL SIGTERM

USER 65532:65532
ENTRYPOINT ["/app/main"]
