# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Huobao Drama is an AI-powered short drama production platform. It automates the entire workflow from script generation, character design, storyboarding to video composition. The architecture follows Domain-Driven Design (DDD) principles with clean separation between layers.

**Tech Stack:**
- **Backend:** Go 1.23+, Gin web framework, GORM ORM, SQLite (with modernc.org/sqlite for CGO-free builds)
- **Frontend:** Vue 3 + TypeScript + Vite, Element Plus UI, TailwindCSS, Pinia state management
- **AI Integration:** OpenAI, Doubao, and other text/image/video generation services

## Development Commands

### Backend Development

```bash
# Run backend directly (development mode with auto-reload if using air)
go run main.go

# Build for production
go build -o huobao-drama .

# Run database migrations (handled automatically on startup via AutoMigrate)
# Manual migration tool available:
go run cmd/migrate/main.go

# Run tests
go test ./...

# Format code
gofmt -w .

# Run linter (if golangci-lint is installed)
golangci-lint run
```

### Frontend Development

```bash
cd web

# Install dependencies
npm install

# Development server (with hot reload at http://localhost:3012)
npm run dev

# Build for production
npm run build

# Type-check without building
npm run build:check

# Lint code
npm run lint
```

### Full Stack Development (Recommended)

Run both services simultaneously:
- Terminal 1: `go run main.go` (backend on port 5678)
- Terminal 2: `cd web && npm run dev` (frontend on port 3012)

The frontend Vite dev server proxies `/api` and `/static` requests to the backend.

### Docker Development

```bash
# Build and start all services
docker compose build
docker compose up -d

# View logs
docker compose logs -f

# Stop services
docker compose down
```

## Architecture Overview

The project follows a **4-layer clean architecture** pattern:

```
api/                    # API Layer - HTTP handlers, middleware, routes
├── handlers/          # Request/response handling
├── middleware/        # CORS, logging, auth
└── routes/           # Route registration

application/           # Application Service Layer - Business orchestration
└── services/         # Use case coordination

domain/               # Domain Layer - Core business logic
├── models/          # Entities and aggregates (Drama, Episode, Character, Scene, etc.)
└── repositories/    # Repository interfaces

infrastructure/       # Infrastructure Layer - External concerns
├── database/        # GORM implementation, migrations
├── storage/         # Local file storage
└── external/        # AI service adapters

pkg/                  # Shared utilities
├── config/          # Configuration loading (Viper)
├── logger/          # Logging (Zap)
└── utils/           # Video/image processing helpers
```

**Key architectural principles:**
- Repository interfaces are defined in `domain/`, implemented in `infrastructure/`
- Domain models are plain Go structs with GORM tags
- Services orchestrate business logic but don't contain domain rules
- AI providers are pluggable through adapter pattern

## Core Domain Models

The system centers around these key entities:

- **Drama:** Top-level project container
- **Episode:** Individual episodes within a drama
- **Character:** Character designs with AI-generated images
- **Scene:** Visual scenes with backgrounds
- **Storyboard:** Frame-by-frame visual planning
- **Asset:** Centralized media storage

## Configuration

Configuration is loaded from `configs/config.yaml` (copy from `config.example.yaml`):

```yaml
app:
  name: "Huobao Drama API"
  debug: true              # Set false for production
  language: "zh"           # zh or en

server:
  port: 5678
  cors_origins:
    - "http://localhost:3012"

database:
  type: "sqlite"
  path: "./data/drama_generator.db"

storage:
  type: "local"
  local_path: "./data/storage"
  base_url: "http://localhost:5678/static"

ai:
  default_text_provider: "openai"
  default_image_provider: "openai"
  default_video_provider: "doubao"
```

**Important:** AI API keys are configured via the web UI, not in config files.

## Database

- Uses SQLite with `modernc.org/sqlite` driver (CGO-free for cross-compilation)
- GORM AutoMigrate runs on startup - no manual migration needed for development
- Connection pooling and WAL mode enabled for concurrency
- Database file: `./data/drama_generator.db`

## Adding New Features

When implementing new features, follow this order:

1. **Define domain model** in `domain/models/`
2. **Create repository interface** in `domain/repositories/`
3. **Implement repository** in `infrastructure/database/`
4. **Create application service** in `application/services/`
5. **Add API handler** in `api/handlers/`
6. **Register routes** in `api/routes/routes.go`

Example pattern for a new entity:
```go
// domain/models/myfeature.go
type MyFeature struct {
    ID        string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
    Name      string    `json:"name" gorm:"type:varchar(100);not null"`
    CreatedAt time.Time `json:"created_at"`
}
```

## Common Issues

- **"database is locked"**: WAL mode should be enabled. Check SQLite file permissions.
- **CORS errors**: Verify `server.cors_origins` in config includes your frontend URL
- **FFmpeg not found**: Install FFmpeg and ensure it's in PATH (required for video processing)
- **Frontend can't reach backend**: In dev mode, Vite proxies requests - check `vite.config.ts`

## Testing

- Backend uses Go's `testing` package
- Frontend uses Vitest (check `web/package.json` for test scripts)
- Test coverage is minimal - add tests when modifying core business logic

## Deployment

Production builds embed the frontend in the Go binary:

```bash
# 1. Build frontend
cd web && npm run build && cd ..

# 2. Build backend (includes web/dist/)
go build -o huobao-drama .

# 3. Run single binary
./huobao-drama
```

The backend serves both API (`/api/v1`) and frontend (`/`) from the same port.
