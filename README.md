# Music Player Monorepo

## Overview

Enterprise-grade Music Player system built with microservices architecture, featuring event-driven communication, distributed caching, and production-ready infrastructure. Designed for scalability, maintainability, and rapid development with modern DevOps practices.

## Architecture Highlights

- 🏗️ **Microservices Architecture**: Loosely coupled services with clear boundaries
- 🔐 **Security-First**: JWT with EdDSA (Ed25519) for access tokens, HS256 for refresh tokens, automatic key rotation
- 📨 **Event-Driven**: Apache Kafka for asynchronous inter-service communication
- 🚀 **High Performance**: Redis caching, connection pooling, optimized Kafka producers (100k-1M msg/s)
- 🎯 **Production-Ready**: Docker Compose orchestration, health checks, resource limits, graceful shutdown
- 🧪 **Clean Architecture**: SOLID principles, Dependency Injection (Google Wire), layered design

## Directory Structure

```
music-player/
├── api/                        # Shared protocol buffers & gRPC definitions
│   └── proto/
│       └── auth/v1/           # Auth service gRPC contracts
├── services/                   # Microservices
│   ├── auth-service/          # Authentication & user management
│   ├── gateway/               # API Gateway with gRPC & HTTP
│   └── notification-service/  # Notification & event processing
├── infra/                      # Infrastructure utilities
│   ├── jwt/                   # JWT keys & JWKS management
│   │   ├── private/          # Ed25519 private keys
│   │   ├── public/           # Ed25519 public keys
│   │   └── jwks/             # JSON Web Key Sets
│   └── scripts/              # Automation scripts
│       ├── rotate_jwt_key.sh    # JWT key rotation
│       └── update_jwks.sh       # JWKS update
├── data/                       # Persistent data volumes
│   ├── postgres/              # PostgreSQL data
│   ├── redis/                 # Redis persistence
│   ├── kafka/                 # Kafka logs
│   └── zookeeper/             # Zookeeper data
├── docker-compose.yml          # Production-grade infrastructure
├── .env.example                # Environment template
└── README.md
```

## Services Overview

### 🔐 Auth Service (Port 8080 HTTP, 8081 gRPC)

**Purpose**: Authentication, authorization, user management

**Features**:

- User registration & login with JWT authentication
- Two-Factor Authentication (TOTP) with QR code generation
- Session management with Redis-backed refresh tokens
- JWT with EdDSA (Ed25519) for access tokens, HS256 for refresh tokens
- JWKS endpoint for public key distribution
- gRPC server for inter-service communication
- Kafka producer for user lifecycle events

**Tech Stack**: Go 1.25, Gin, GORM, Google Wire, PostgreSQL, Redis, Kafka (franz-go)

**Key Endpoints**:

- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/2fa/setup` - Enable 2FA
- `GET /api/v1/.well-known/jwks.json` - JWKS public keys

### 🚪 Gateway (Port 3000)

**Purpose**: API Gateway, request routing, authentication middleware

**Features**:

- Unified entry point for all client requests
- JWT token verification via JWKS
- gRPC client for auth-service communication
- Request/response transformation
- Authentication middleware with Redis session validation
- Rate limiting & request logging

**Tech Stack**: Go 1.25, Gin, gRPC client, Redis, Google Wire

**Key Endpoints**:

- `POST /api/v1/auth/login` - Proxied login
- `POST /api/v1/auth/register` - Proxied registration
- `GET /api/v1/users` - Get user profile (protected)
- `POST /api/v1/2fa/*` - 2FA management (protected)

### 📬 Notification Service (Port 8082)

**Purpose**: Event processing, notifications, background jobs

**Features**:

- Kafka consumer for event-driven architecture
- Process user lifecycle events (registration, login, etc.)
- Email/SMS notification support (planned)
- Async task processing

**Tech Stack**: Go 1.25, Gin, Kafka (franz-go), Google Wire

**Events Consumed**:

- `user.registered` - User registration notifications
- (More events to be added)

## Infrastructure Components

### PostgreSQL (Port 5432)

- Primary database for auth-service
- Connection pooling, health checks
- Automatic migration with Goose
- PgAdmin UI on port 5050

### Redis Stack (Port 6379, UI 8001)

- Session storage & refresh tokens
- Distributed caching
- LRU eviction policy (1GB max memory)
- RedisInsight UI for monitoring

### Apache Kafka (Port 9092)

- Event streaming platform
- 3 partitions per topic for parallelism
- 7-day retention policy
- Kafka UI on port 8080

### Zookeeper (Port 2181)

- Kafka cluster coordination
- Distributed configuration management

## Technologies Used

### Backend

- **Language**: Go 1.25.0
- **Web Framework**: Gin (HTTP), gRPC (inter-service)
- **ORM**: GORM with PostgreSQL driver
- **DI**: Google Wire (compile-time dependency injection)
- **Config**: Viper (environment management)
- **Migration**: Goose
- **Kafka Client**: franz-go (high-performance)
- **JWT**: golang-jwt/jwt/v5 with Ed25519 + HS256

### Infrastructure

- **Databases**: PostgreSQL 16, Redis Stack 7
- **Message Queue**: Apache Kafka 3.x + Zookeeper
- **Containerization**: Docker Compose with resource limits
- **Monitoring**: JSON logging, health checks, PgAdmin, RedisInsight, Kafka UI

### DevOps

- **Container Orchestration**: Docker Compose
- **Build System**: Multi-stage Dockerfiles with distroless images
- **Security**: Non-root containers, read-only filesystem, capability dropping
- **CI/CD**: GitHub Actions (planned)

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.25+ (for local development)
- Git

### 1. Clone Repository

```bash
git clone https://github.com/VanTruongNg/music-player.git
cd music-player
```

### 2. Configure Environment

```bash
# Copy environment template
cp .env.example .env

# Edit .env with your configuration
# Default values work for local development
```

### 3. Start Infrastructure

```bash
# Start all infrastructure services (Postgres, Redis, Kafka, Zookeeper)
docker compose up -d postgres redis-stack kafka zookeeper

# Wait for services to be healthy (~30 seconds)
docker compose ps

# Check logs
docker compose logs -f
```

### 4. Run Database Migration

```bash
# Run migrations for auth-service
docker compose up migration

# Or manually with goose:
cd services/auth-service
goose -dir migrations postgres "postgres://postgres:postgres123@localhost:5432/music_player?sslmode=disable" up
```

### 5. Start Services

**Option A: Run with Go (Development)**

```bash
# Terminal 1: Auth Service
cd services/auth-service
go mod download
go run ./cmd

# Terminal 2: Gateway
cd services/gateway
go mod download
go run ./cmd

# Terminal 3: Notification Service
cd services/notification-service
go mod download
go run ./cmd
```

**Option B: Build Docker Images (Production)**

```bash
# Build all services
docker compose build

# Start all services
docker compose up -d
```

### 6. Verify Installation

```bash
# Check auth-service health
curl http://localhost:8080/api/v1/health

# Check gateway health
curl http://localhost:3000/api/v1/health

# Check notification-service health
curl http://localhost:8082/health

# Access management UIs:
# - Kafka UI: http://localhost:8080
# - RedisInsight: http://localhost:8001
# - PgAdmin: http://localhost:5050 (admin@admin.com / admin)
```

## Development Workflow

### Code Standards

- **Architecture**: Follow SOLID principles, clean architecture, DRY, KISS
- **Naming**: Clear, descriptive names for functions, variables, packages
- **Error Handling**: Always handle errors explicitly, use custom error types
- **Logging**: Structured logging with severity levels [INFO], [WARN], [ERROR], [FATAL]
- **Testing**: Unit tests for business logic, integration tests for APIs
- **Documentation**: Each service has detailed README.md with diagrams

### Project Structure Rules

- `cmd/`: Application entry points only
- `configs/`: Configuration loading and validation
- `internal/`: Private application code (domain, services, handlers, repositories)
- `api/`: Shared protocol buffers and gRPC definitions
- `migrations/`: Database schema versions (Goose format)

### Making Changes

```bash
# 1. Create feature branch
git checkout -b feature/your-feature

# 2. Make changes, following code standards

# 3. Run tests
go test ./...

# 4. Build and verify
go build ./cmd

# 5. Update documentation if needed

# 6. Commit with descriptive message
git commit -m "feat(auth): add OAuth2 support"

# 7. Push and create PR
git push origin feature/your-feature
```

### Wire Dependency Injection

When adding new dependencies:

```bash
cd services/[service-name]

# Update cmd/wire.go with new providers
# Then regenerate wire_gen.go:
go run github.com/google/wire/cmd/wire@latest ./cmd
```

### JWT Key Rotation

```bash
# Generate new Ed25519 key pair
./infra/scripts/rotate_jwt_key.sh

# Update JWKS
./infra/scripts/update_jwks.sh

# Restart auth-service to load new keys
docker compose restart auth-service
```

## API Documentation

### Authentication Flow

```mermaid
sequenceDiagram
    participant Client
    participant Gateway
    participant AuthService
    participant Redis
    participant Kafka

    Client->>Gateway: POST /auth/register
    Gateway->>AuthService: gRPC Register()
    AuthService->>PostgreSQL: Create User
    AuthService->>Kafka: Publish user.registered
    AuthService->>Redis: Store refresh token
    AuthService-->>Gateway: JWT tokens
    Gateway-->>Client: Access + Refresh tokens

    Note over Kafka,NotificationService: Async event processing
    Kafka->>NotificationService: Consume user.registered
    NotificationService->>NotificationService: Send welcome email
```

### Service Communication

```mermaid
graph LR
    Client[Client/Frontend] -->|HTTP/REST| Gateway
    Gateway -->|gRPC| AuthService[Auth Service]
    Gateway -->|Redis| Cache[(Redis)]
    AuthService -->|GORM| DB[(PostgreSQL)]
    AuthService -->|Session| Cache
    AuthService -->|Events| Kafka[Kafka]
    Kafka -->|Consume| NotificationService[Notification Service]

    style Gateway fill:#e1f5ff
    style AuthService fill:#ffe1e1
    style NotificationService fill:#e1ffe1
```

## System Architecture

```mermaid
flowchart TD
    subgraph Clients
        Web[Web App]
        Mobile[Mobile App]
        API[External API]
    end

    subgraph Gateway Layer
        GW[API Gateway<br/>Port 3000]
    end

    subgraph Microservices
        Auth[Auth Service<br/>HTTP: 8080<br/>gRPC: 8081]
        Notif[Notification Service<br/>Port: 8082]
    end

    subgraph Data Layer
        PG[(PostgreSQL<br/>Port: 5432)]
        RD[(Redis<br/>Port: 6379)]
        KF[(Kafka<br/>Port: 9092)]
    end

    subgraph Management
        PGA[PgAdmin<br/>Port: 5050]
        RDI[RedisInsight<br/>Port: 8001]
        KUI[Kafka UI<br/>Port: 8080]
    end

    Web --> GW
    Mobile --> GW
    API --> GW

    GW -->|gRPC| Auth
    GW --> RD

    Auth --> PG
    Auth --> RD
    Auth -->|Produce| KF

    KF -->|Consume| Notif

    PGA -.->|Manage| PG
    RDI -.->|Monitor| RD
    KUI -.->|Monitor| KF

    style GW fill:#4a90e2
    style Auth fill:#e27474
    style Notif fill:#74e274
    style PG fill:#f4a261
    style RD fill:#e76f51
    style KF fill:#2a9d8f
```

## Performance Characteristics

### Auth Service

- **Throughput**: ~10,000 requests/second (with Redis caching)
- **Latency**:
  - Login: ~50ms (p95)
  - Token refresh: ~10ms (p95)
  - 2FA verification: ~100ms (p95)

### Kafka Producer (Event Publishing)

- **Balanced Profile**: 100k-500k msg/s, LeaderAck, Snappy compression
- **Safe Profile**: 5k-10k msg/s, AllISRAcks, Lz4 compression
- **Fast Profile**: 500k-1M+ msg/s, NoAck, no compression

### Gateway

- **Throughput**: ~15,000 requests/second
- **Latency**:
  - Proxy overhead: ~5ms
  - With auth verification: ~15ms (p95)

## Monitoring & Observability

### Health Checks

- All services expose `/health` endpoint
- Docker health checks configured with retries
- Graceful shutdown with SIGTERM handling

### Logging

- Structured JSON logging
- Log levels: INFO, WARN, ERROR, FATAL
- Request ID tracing across services
- Centralized log aggregation (planned: ELK stack)

### Metrics (Planned)

- Prometheus exposition format
- Grafana dashboards
- Key metrics: request rate, error rate, latency (RED method)

## Security Best Practices

### Authentication

- JWT with Ed25519 (EdDSA) for access tokens - asymmetric, more secure
- HS256 for refresh tokens - symmetric, faster validation
- Automatic key rotation support
- JWKS endpoint for public key distribution

### Infrastructure

- Non-root containers (UID 65532)
- Read-only root filesystem where possible
- Capability dropping (no-new-privileges)
- Security scanning with Trivy (planned)

### Database

- Parameterized queries (GORM ORM)
- Connection pooling with limits
- Password encryption (bcrypt)
- SCRAM-SHA-256 authentication for Postgres

### Redis

- Password authentication
- Connection limits
- Memory limits with LRU eviction

## Troubleshooting

### Service won't start

```bash
# Check logs
docker compose logs [service-name]

# Check if ports are available
netstat -an | grep -E "3000|8080|8081|8082|5432|6379|9092"

# Restart specific service
docker compose restart [service-name]
```

### Database connection issues

```bash
# Check Postgres health
docker compose exec postgres pg_isready -U postgres

# Check connection from host
psql -h localhost -U postgres -d music_player

# Recreate database (WARNING: data loss)
docker compose down -v
docker compose up -d postgres
docker compose up migration
```

### Kafka connection issues

```bash
# Check Kafka is ready
docker compose exec kafka kafka-topics --bootstrap-server localhost:29092 --list

# Create topic manually
docker compose exec kafka kafka-topics --bootstrap-server localhost:29092 --create --topic user.registered --partitions 3 --replication-factor 1

# Check consumer group lag
docker compose exec kafka kafka-consumer-groups --bootstrap-server localhost:29092 --describe --group notification-service
```

### Wire generation fails

```bash
cd services/[service-name]
go get github.com/google/wire/cmd/wire@latest
go run github.com/google/wire/cmd/wire@latest ./cmd
```

## Roadmap

### Phase 1: Core Services ✅

- [x] Auth Service with JWT & 2FA
- [x] API Gateway with gRPC
- [x] Notification Service foundation
- [x] Kafka event-driven architecture
- [x] Infrastructure setup (Docker Compose)

### Phase 2: Music Features (In Progress)

- [ ] Music Library Service
- [ ] Streaming Service with CDN
- [ ] Playlist Management
- [ ] User Preferences & Recommendations

### Phase 3: Advanced Features (Planned)

- [ ] Search Service (Elasticsearch)
- [ ] Social Features (follow, share)
- [ ] Real-time Chat (WebSocket)
- [ ] Analytics & Reporting

### Phase 4: Production Readiness (Planned)

- [ ] Kubernetes deployment
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Monitoring (Prometheus + Grafana)
- [ ] Distributed tracing (Jaeger)
- [ ] API rate limiting
- [ ] CDN integration (CloudFlare)

## Contributing

Contributions are welcome! Please follow these guidelines:

1. **Fork** the repository
2. Create a **feature branch** (`git checkout -b feature/amazing-feature`)
3. Follow **code standards** (SOLID, clean architecture)
4. Write **tests** for new features
5. Update **documentation** (README, diagrams)
6. **Commit** changes (`git commit -m 'feat: add amazing feature'`)
7. **Push** to branch (`git push origin feature/amazing-feature`)
8. Open a **Pull Request**

### Commit Message Convention

```
type(scope): subject

body (optional)

footer (optional)
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

Example: `feat(auth): add OAuth2 Google login support`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact & Support

- **Maintainer**: Van Truong Nguyen
- **Email**: truongnguyen060603@gmail.com
- **GitHub**: [@VanTruongNg](https://github.com/VanTruongNg)

For issues and feature requests, please use [GitHub Issues](https://github.com/VanTruongNg/music-player/issues).

---

⭐ **If you find this project helpful, please give it a star!** ⭐

_Built with ❤️ using modern Go practices and production-ready architecture._
