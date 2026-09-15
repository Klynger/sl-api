# Shopping list API

A REST API for a shopping list application, written in Go with chi, GORM, and PostgreSQL.

## Prerequisites

- **Go** 1.21 or newer: any recent release works, since the `go` command automatically downloads the exact toolchain version pinned in `go.mod`
- **Docker** with **Compose v2** (the `docker compose` subcommand), used for Postgres and the full-stack run
- **curl**, used by the verification script's smoke test

Check everything at once with:

```bash
./scripts/check-prereqs.sh
```

## Quick start

```bash
# Run everything (Postgres + migrations + API)
docker compose up --build

# Or build and test locally
go build -o ./bin/api ./cmd/api
go test ./...

# Full verification suite (what CI runs): lint, build, tests,
# migration cycle, and an API smoke test against a fresh Postgres
./scripts/verify.sh
```

Configuration is environment based; see `config/config.go` for the required variables.

## Documentation

- [Architecture](docs/architecture.md): high-level diagram and layer responsibilities
- [Current state](docs/current-state.md): implemented features, endpoints, and known gaps
