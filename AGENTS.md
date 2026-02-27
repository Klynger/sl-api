# AGENTS.md — sl-api

Go REST API for a shopping list application, built with chi router, GORM, and PostgreSQL.

## Build & Run

```bash
# Build the API server
go build -o ./bin/api ./cmd/api

# Build the migration CLI
go build -o ./bin/migrate ./cmd/migrate

# Run with Docker Compose (starts Postgres + runs migrations + starts API)
docker compose up --build
```

## Test Commands

```bash
# Run all tests
go test ./...

# Run a single test by name
go test ./api/resource/product/ -run TestRepository_Create

# Run tests in a specific package
go test ./api/resource/product/

# Verbose output
go test -v ./...

# Run tests with race detector
go test -race ./...
```

There is currently one test file: `api/resource/product/repository_test.go`. Tests use
the standard `testing` package with custom helpers from `util/test/` and DB mocks from
`mock/db/` (go-sqlmock). Test files use external test packages (e.g., `package product_test`).

## Linting & Formatting

```bash
# Format Go code (standard)
gofmt -w .

# Vet for common issues
go vet ./...

# SQL linting (config in .sqlfluff, postgres dialect)
sqlfluff lint migrations/
```

No golangci-lint configuration exists. Rely on `gofmt` and `go vet`.

## Project Structure

```
cmd/
  api/main.go            # HTTP server entrypoint
  migrate/main.go        # Migration CLI entrypoint
config/
  config.go              # Env-based configuration (joeshaw/envdecode)
api/
  middleware/auth.go      # Session-based auth middleware
  resource/
    common/err/err.go     # Shared HTTP error response helpers
    health/handler.go     # Health check endpoint (GET /livez)
    product/              # Product domain: handler, model, repository, tests
    user/                 # User domain: handler, model, repository
mock/db/db.go            # GORM mock DB helper (go-sqlmock)
util/
  test/test.go           # Custom test assertions (NoError, Equal)
  validator/             # go-playground/validator setup + error formatting
migrations/              # Goose SQL migration files
```

## Code Style

### Import Organization

Group imports in three blocks separated by blank lines: stdlib, third-party, local module.

```go
import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"
    "gorm.io/gorm"

    e "sl-api/api/resource/common/err"
    "sl-api/util/validator"
)
```

Use short aliases when needed to avoid conflicts:
- `e` for `sl-api/api/resource/common/err`
- `validatorUtil` for `sl-api/util/validator`
- `mockDB` for `sl-api/mock/db`
- `testUtil` for `sl-api/util/test`

### Naming Conventions

- **Packages:** lowercase, short (`product`, `user`, `health`, `err`, `config`)
- **Structs:** PascalCase (`Product`, `Repository`, `API`)
- **Constructors:** `New()` or `NewRepository(db)` pattern
- **Handler methods:** Named by action: `List`, `Read`, `Create`, `Update`, `Delete`, `Register`, `Login`, `Logout`
- **Receiver names:** Single letter matching type (`r` for Repository, `a` for API, `p` for Product)
- **Config fields:** Use `env:"FIELD_NAME,required"` struct tags

### Architecture Pattern — Resource-Based

Each domain entity gets its own package under `api/resource/` with:

| File | Purpose |
|------|---------|
| `handler.go` | HTTP handler — `API` struct with route handler methods |
| `model.go` | GORM model, DTO (output), Form (input), conversion methods |
| `repository.go` | Data access layer wrapping GORM |
| `repository_test.go` | Tests for repository using go-sqlmock |

The `API` struct holds a `*Repository` and `*validator.Validate`, constructed via `New(db, v)`.

### Model / DTO / Form Convention

- **Model:** GORM struct with DB tags, timestamps, `gorm.DeletedAt` for soft deletes
- **DTO:** Output struct with `json` tags (omit sensitive fields like passwords)
- **Form:** Input struct with `json` + `validate` tags for request binding
- **Conversions:** `model.ToDto()` and `form.ToModel()` methods
- **Collections:** Define named slice types (e.g., `type Products []*Product`) with their own `ToDto()` method

### Error Handling

- Use centralized helpers from `api/resource/common/err`:
  - `e.ServerError(w, msg)` — 500 response
  - `e.BadRequest(w, msg)` — 400 response
  - `e.ValidationErrors(w, msg)` — 422 response
- Pre-defined byte slices for common error messages (`RespDBDataAccessFailure`, `RespJSONEncodeFailure`, etc.)
- Pattern in handlers: check error, write error response, `return` early
- Fatal errors during startup use `log.Fatal()` / `log.Fatalf()`

### Testing Conventions

- Use external test packages (`package product_test`, not `package product`)
- Use `t.Parallel()` for test functions
- Use custom assertions from `util/test/`: `testUtil.NoError(t, err)`, `testUtil.Equal(t, x, y)`
- Mock GORM with `mock/db/db.go`: `mockDB.NewMockDB()` returns `(*gorm.DB, sqlmock.Sqlmock)`
- Match SQL with `regexp.QuoteMeta()` when setting mock expectations

## Key Dependencies

`go-chi/chi/v5` (router), `gorm.io/gorm` + `gorm.io/driver/postgres` (ORM),
`pressly/goose/v3` (migrations), `go-playground/validator/v10` (validation),
`google/uuid`, `gorilla/sessions` (auth), `joeshaw/envdecode` (config),
`DATA-DOG/go-sqlmock` (test mocking).
