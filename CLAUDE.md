# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

API service for managing stray dog trapping requests (отлов бродячих собак). Built with Go (Gin framework) and a Python worker microservice for Excel document generation.

## Architecture

The system consists of two services orchestrated via docker-compose:

- **app** (Go, port 8091) — main REST API handling auth, CRUD for requests, users, districts, territorial departments (ter_otdels), applicants, and roles
- **worker** (Python/FastAPI, port 8001) — generates Excel documents from templates using openpyxl, communicates via HTTP from the Go service

### Go service layering (internal/):
- `application/` — HTTP handlers (Gin), middleware, router setup
- `service/` — business logic, depends on repository and storage interfaces
- `repository/` — PostgreSQL queries using Masterminds/squirrel query builder
- `dto/` — data transfer objects shared across layers
- `config/` — env-based configuration (godotenv)
- `s3/` — AWS S3 file storage client
- `worker/` — HTTP client to call the Python worker

### Shared packages (pkg/):
- `jwt/` — RSA-based JWT token service (private/public key from `certs/`)
- `postgres/` — pgxpool connection setup
- `security/` — password hashing (bcrypt)
- `validator/` — request validation wrapper

## Build & Run

```bash
# Build the Go binary
go build -o app_service cmd/main.go

# Run with Docker Compose (starts both app and worker)
docker-compose up --build

# Run migrations (standalone migrator tool)
go run migrator.go
```

## Key Technical Details

- Database: PostgreSQL, connected via pgx/v5 pool
- Query builder: Masterminds/squirrel with `PlaceholderFormat(sq.Dollar)`
- Auth: RSA JWT (keys in `certs/private.pem` and `certs/public.pem`)
- Config: environment variables loaded from `.env` via godotenv
- API docs: Swagger annotations (swaggo/swag), generated files in `docs/`
- Migrations: golang-migrate/v4, SQL files in `migrations/`
- The Go and Python services share files via a Docker volume (`shared-data` mounted at `/app/shared`)

## Swagger Docs

```bash
# Regenerate swagger docs after changing annotations
swag init -g cmd/main.go
```

## Database Schema

Main tables: `roles`, `districts`, `ter_otdels`, `users`, `applicants`, `external_users`, `requests`. See `migrations/000001_init.up.sql` for full schema. All primary keys are UUIDs (gen_random_uuid()).
