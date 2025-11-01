# Auth Service

## Overview

**Auth Service** is a robust user authentication microservice built with Golang, featuring modular architecture, multi-factor authentication (2FA), JWT with EdDSA signing for access tokens and HS256 for refresh tokens, session management, and integration with PostgreSQL, Redis, and Kafka. This service is the security foundation for the Music Player system, ensuring safety, scalability, and maintainability.

## Key Components

- **cmd/**: Application entry point, initializes configuration, router, DI.
- **configs/**: Manages app, database, redis, kafka configuration.
- **internal/domain/**: Defines entities, errors, core business logic.
- **internal/dto/**: Defines API request/response structs.
- **internal/handlers/**: Handles HTTP requests for user, 2FA.
- **internal/repositories/**: Data access (PostgreSQL) for users.
- **internal/services/**: Business logic for user, 2FA, token management.
- **internal/utils/jwt/**: JWT utilities with EdDSA for access tokens, HS256 for refresh tokens, key rotation support.
- **migrations/**: Database schema and migration scripts.
- **docker-compose.yml**: Quick start for Redis, Postgres, Kafka, Zookeeper, Kafka UI, PgAdmin.

## Quick Start

# Auth Service

## Overview

Auth Service is the authentication and user management microservice for the Music Player system. It provides user registration, login, JWT issuance, 2FA (TOTP), session management (Redis), and publishes user lifecycle events to Kafka.

This README documents the current, implemented features only (no planned items). For architecture-level docs see the repo-level README.

## What exists in the codebase

- HTTP API implemented with Gin
- gRPC server for inter-service calls (used by the Gateway)
- JWT signing using Ed25519 for access tokens and HS256 for refresh tokens
- JWKS endpoint to expose public keys
- 2FA TOTP setup and verification
- Redis-backed refresh token/session management
- Kafka producer integration with an `EventPublisher` service (sync publish supported)
- Google Wire for dependency injection

## Repository layout (relevant folders)

```
services/auth-service/
├── cmd/                  # entry point, wire setup
├── configs/              # app, db, redis, kafka configs
├── internal/
│   ├── handlers/         # HTTP handlers
│   ├── repositories/     # GORM repositories
│   ├── services/         # business logic (UserService, TwoFA, EventPublisher)
│   ├── kafka/            # envelope, producer integration
│   └── utils/            # jwt, redis helpers
├── migrations/           # DB migrations (goose)
├── Dockerfile
└── README.md
```

## Available endpoints (as implemented)

Base path: `/api/v1`

Auth and user routes (in `internal/routes` and `handlers`):

- POST `/api/v1/auth/register` - register a new user (public)
- POST `/api/v1/auth/login` - login and get access+refresh tokens (public)
- POST `/api/v1/auth/refresh` - exchange refresh token for new access token (public)
- POST `/api/v1/auth/logout` - logout, revoke session (protected)
- GET `/api/v1/auth/validate` - validate current token (protected)

2FA management (protected, require auth middleware):

- POST `/api/v1/auth/:id/2fa/setup` - generate TOTP secret / QR data
- POST `/api/v1/auth/:id/2fa/enable` - enable 2FA with OTP
- POST `/api/v1/auth/:id/2fa/verify` - verify OTP
- POST `/api/v1/auth/:id/2fa/disable` - disable 2FA

User management (protected):

- GET `/api/v1/auth/users` - list users (requires appropriate permissions)
- GET `/api/v1/auth/users/:id` - get user by ID
- GET `/api/v1/auth/me` - get current authenticated user

JWKS endpoint:

- GET `/.well-known/jwks.json` - JWKS public keys for verifying access tokens

Notes:

- Routes and handlers are implemented under `internal/routes` and `internal/handlers`.

## Event publishing

The codebase contains an `EventPublisher` service (in `internal/services`) that encapsulates Kafka publishing logic. The `UserService` calls this service after a successful registration to publish a `user.registered` event. The envelope format, serializer, and producer profiles live under `internal/kafka`.

Published topic:

- `user.registered` (JSON envelope)

Producer behavior:

- The project provides sync (`Publish`) and async (`PublishAsync`) publishing APIs. Current production code uses synchronous publish for reliability (errors are handled and logged).

## Configuration

Environment variables are loaded via Viper. Key variables (see `configs/`):

- `APP_PORT` - HTTP port (default: 8080)
- `KAFKA_BROKERS` - comma-separated list of brokers
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`
- `POSTGRES_*` - database connection

Use the top-level `.env.example` as a template.

## Quick start (development)

1. Start infra (Postgres, Redis, Kafka):

```powershell
cd \\path\\to\\music-player
docker compose up -d postgres redis-stack zookeeper kafka
```

2. Prepare .env for auth-service:

```powershell
cd services/auth-service
cp ..\\.env.example .env
# Edit .env if needed
```

3. Run DB migrations (goose):

```powershell
docker compose up migration
# or locally (requires goose installed):
goose -dir migrations postgres "postgres://postgres:postgres123@localhost:5432/music_player?sslmode=disable" up
```

4. Run the service locally:

```powershell
cd services/auth-service
go mod download
go run ./cmd
```

The HTTP server will listen on the port defined in your `.env` (default 8080). The gRPC server listens on a separate port defined in configs (default 8081).

## Development notes

- Dependency injection uses Google Wire. If you modify `cmd/wire.go`, regenerate with:

```powershell
cd services/auth-service
go run github.com/google/wire/cmd/wire@latest ./cmd
```

- The `EventPublisher` is injected into `UserService` via Wire. This makes it easy to swap a mock implementation for unit tests.

- Token signing keys are stored under `infra/jwt/` in the monorepo; key rotation scripts are available in `infra/scripts`.

## Troubleshooting

- `500` on register/login: check Postgres connection and migrations
- Kafka publish warnings: check `KAFKA_BROKERS` and broker health; verify topic exists (`user.registered`)
- JWT verification issues: ensure JWKS endpoint is reachable and keys are present in `infra/jwt/public`

## Where to look in the code

- Routes: `internal/routes`
- Handlers: `internal/handlers`
- Business logic: `internal/services`
- Kafka envelope & producer: `internal/kafka`
- Configs: `configs`

## Contact

- Maintainer: Van Truong Nguyen
- Email: truongnguyen060603@gmail.com

---

_This README documents features present in the code as of this change._
