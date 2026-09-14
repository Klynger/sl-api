# Shopping list API

A REST API for a shopping list application, written in Go with chi, GORM, and PostgreSQL.

## Quick start

```bash
# Run everything (Postgres + migrations + API)
docker compose up --build

# Or build and test locally
go build -o ./bin/api ./cmd/api
go test ./...
```

Configuration is environment based; see `config/config.go` for the required variables.

## Documentation

- [Architecture](docs/architecture.md): high-level diagram and layer responsibilities
- [Current state](docs/current-state.md): implemented features, endpoints, and known gaps
