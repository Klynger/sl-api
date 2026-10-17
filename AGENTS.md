# AGENTS.md — sl-api

Go REST API for a shopping list application, built with chi router, GORM, and PostgreSQL.

## Build & Run

```bash
# Build the API server
go build -o ./bin/api ./cmd/api

# Build the migration CLI
go build -o ./bin/migrate ./cmd/migrate

# Build the local-seed CLI
go build -o ./bin/seed ./cmd/seed

# Run with Docker Compose (starts Postgres + runs migrations + starts API + seeds local data)
docker compose up --build
```

## Test Commands

```bash
# Run all tests
go test ./...

# Run a single test by name
go test ./api/repositories/product/ -run TestRepository_Create

# Run tests in a specific package
go test ./api/services/list/

# Verbose output
go test -v ./...

# Run tests with race detector
go test -race ./...
```

Tests live across the repository and service packages (for example
`api/repositories/product`, `api/repositories/group_member`, `api/services/list`,
`api/services/group`), plus a user model test. They use the standard `testing`
package with custom helpers from `util/test/` and DB mocks from `mock/db/`
(go-sqlmock), in external test packages (e.g., `package product_test`). Handlers
have no unit tests; they are exercised end-to-end against a real Postgres.

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
  api/main.go              # HTTP server entrypoint
  migrate/main.go          # Migration CLI entrypoint
  seed/main.go             # Local-seed CLI (drives the real HTTP endpoints)
config/
  config.go                # Env-based configuration (joeshaw/envdecode)
api/
  middleware/              # Session auth (RequireAuth), default headers, login rate limit
  model/<domain>/          # GORM models, DTOs (output), Forms (input), conversions
  repositories/<domain>/   # Data access, wrapping GORM
  routes/
    common/errors/         # Shared HTTP error response helpers (errorsCommon)
    health/                # Health check endpoint (GET /livez)
    router.go              # Mounts everything under /v1 (public + authenticated groups)
    <domain>/              # HTTP handlers: product, user, group, group/invite, list
  services/
    group/                 # GroupService: groups, members, invites
    list/                  # ListService: lists, items, quantity, status
    util/                  # WithTransaction helper
mock/db/db.go              # GORM mock DB helper (go-sqlmock)
util/
  test/test.go             # Custom test assertions (NoError, Equal)
  validator/               # go-playground/validator setup + error formatting
migrations/                # Goose SQL migration files
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

    errorsCommon "sl-api/api/routes/common/errors"
    validatorUtil "sl-api/util/validator"
)
```

Use short aliases when needed to avoid conflicts:
- `errorsCommon` for `sl-api/api/routes/common/errors`
- `validatorUtil` for `sl-api/util/validator`
- `mockDB` for `sl-api/mock/db`
- `testUtil` for `sl-api/util/test`
- `<domain>Service` for `sl-api/api/services/<domain>` (e.g. `listService`, `groupService`), `<domain>Repository` for repositories, `<domain>Model` for models

### Naming Conventions

- **Packages:** lowercase; handler packages use a `<domain>Handlers` name, services `<domain>Service`, repositories `<domain>Repository`, models `<domain>Model`
- **Structs:** PascalCase (`Product`, `Repository`, `API`)
- **Constructors:** `New()` or `NewRepository(db)` pattern
- **Handler methods:** Named by action: `List`, `Read`, `Create`, `Update`, `Delete`, `Register`, `Login`, `Logout`
- **Receiver names:** Single letter matching type (`r` for Repository, `a` for API, `p` for Product)
- **Config fields:** Use `env:"FIELD_NAME,required"` struct tags

### Architecture Pattern — Layered

A domain is split across four sibling trees rather than one package:

| Location | Purpose |
|----------|---------|
| `api/model/<domain>/` | GORM model, DTO (output), Form (input), conversion methods |
| `api/repositories/<domain>/` | Data access wrapping GORM, plus `_test.go` (go-sqlmock) |
| `api/routes/<domain>/` | HTTP handlers — an `API` struct with route handler methods |
| `api/services/<domain>/` | Business logic + authorization (only for richer domains) |

The handler `API` struct holds its dependencies (a repository and/or a service, plus `*validator.Validate`), constructed via `New(...)`. Simple domains (products, users) skip the service layer and call repositories directly; groups/invites (`GroupService`) and lists/items (`ListService`) have a service.

### Service Layer Structure

Each domain has exactly one service struct under `api/services/<domain>/`
(`ListService`, `GroupService`), constructed with `New(db)`. Every operation is a
method on that struct, one method per file (`create_list.go`, `add_item.go`,
`list_group_invites.go`, ...), named by its action (`CreateList`, `AddItem`,
`ListGroupInvites`). The operation's input and output types live in the same file
as the method and are named for the operation (`CreateListInput`,
`CreateInviteOutput`). Shared helpers live once on the struct, e.g.
`ensureGroupMember` for the group-membership check.

Do **not** add a separate struct-plus-`Execute` per operation; the single
service-per-domain shape is the standard (decided in KLY-19). Handlers hold the
domain's service as a single field, wired through its `New(db)`.

### Model / DTO / Form Convention

- **Model:** GORM struct with DB tags, timestamps, `gorm.DeletedAt` for soft deletes (the one exception is `list_items`, deliberately hard-deleted, so it has no `DeletedAt`)
- **DTO:** Output struct with `json` tags (omit sensitive fields like passwords)
- **Form:** Input struct with `json` + `validate` tags for request binding
- **Conversions:** `model.ToDto()` and `form.ToModel()` methods
- **Collections:** Define named slice types (e.g., `type Products []*Product`) with their own `ToDto()` method

### Error Handling

- Use centralized helpers from `api/routes/common/errors` (aliased `errorsCommon`):
  - `errorsCommon.ServerError(w, item)` — 500, `BadRequest` — 400, `Unauthorized` — 401, `Forbidden` — 403, `NotFound` — 404, `Conflict` — 409, `ValidationErrors` — 422.
  - Build items with `errorsCommon.NewError(code, detail)` / `NewErrorWithMeta(code, detail, meta)`; `code` constants live in that package. The wire shape is `{"errors": [{code, detail, meta}]}`.
  - `errorsCommon.WriteDBError(w, err, fallback)` classifies Postgres constraint violations (duplicate → 409, FK violation → 422) and otherwise writes the fallback as 500.
- Services return sentinel errors (`errors.New`, one per package `errors.go`); handlers map them to responses with a `switch { case errors.Is(err, ...) }`.
- Pattern in handlers: check error, write error response, `return` early.
- Fatal errors during startup use `log.Fatal()` / `log.Fatalf()`.

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
`jackc/pgx/v5` (Postgres driver / error classification via `pgconn`),
`golang.org/x/sync/errgroup` (concurrent seed), `DATA-DOG/go-sqlmock` (test mocking).
