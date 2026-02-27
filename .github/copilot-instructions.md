# Copilot Instructions

## Commands

```bash
make run          # Run the application
make build        # Build binary to bin/hyperscal-go
make test         # Run all tests
make dev          # Hot reload (requires air)
make install      # Download and tidy dependencies
make docker-postgres  # Start local PostgreSQL via Docker
```

Run a single test:
```bash
go test -v -run TestFunctionName ./internal/service/
```

## Architecture

This is a layered REST API (Gin + GORM) with **driver-switchable database support** (PostgreSQL / Oracle).

```
main.go              → wires everything: config → DB → repo → service → controller → router
config/              → loads all env vars into typed Config struct
pkg/database/        → connects DB and runs GORM AutoMigrate on startup
internal/
  domain/            → GORM model structs (Country, City, User); each has TableName()
  repository/        → interfaces (e.g. CountryRepository); implementations under postgres/ and oracle/
  service/           → business logic; converts between domain and dto
  controller/        → Gin handlers; binds request DTOs, calls service, returns APIResponse
  dto/               → request/response structs; APIResponse helpers (SuccessResponse / ErrorResponse)
pkg/
  jwt/               → GenerateToken / ValidateToken using golang-jwt/jwt v5
  hash/              → bcrypt helpers
  middleware/        → JWTAuthMiddleware (sets user_id and user_email in gin.Context)
```

The DB driver is selected at startup via `DB_DRIVER` env var. When adding a new entity, provide both a `postgres/` and `oracle/` implementation of the repository interface.

## Key Conventions

**Response shape** — all endpoints return `dto.APIResponse`:
```json
{ "success": true/false, "message": "...", "data": ..., "error": "..." }
```
Use `dto.SuccessResponse(message, data)` and `dto.ErrorResponse(message, errString)`.

**Route groups** — public routes live under `/api/auth`; everything else is under a `protected` group that applies `JWTAuthMiddleware()`. Authenticated handlers can read `ctx.GetString("user_id")` and `ctx.GetString("user_email")`.

**Domain models** — always implement `TableName() string` to explicitly set the table name.

**Service → Repository** — services accept/return DTOs, never raw domain models to callers. The private `toResponse()` method on each service handles the conversion.

**Pagination** — use `dto.PaginationRequest` (embedded) in search request DTOs and return `dto.PaginationResponse`.

**Configuration** — copy `.env.example` to `.env`. All config is loaded once via `config.LoadConfig()`; the `pkg/jwt` package calls `LoadConfig()` internally on each token operation.

**Database migrations** — GORM `AutoMigrate` runs automatically on startup. Add new domain models to the `autoMigrate()` call in `pkg/database/database.go`.

**Oracle support** — currently stubbed; requires Oracle Instant Client and the `godror` driver. Only PostgreSQL is functional.
