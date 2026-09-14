#!/usr/bin/env bash
# Single source of truth for CI. Runs the full verification suite:
# formatting, vet, build, unit tests, the migration cycle against a fresh
# Postgres, and an end-to-end smoke test of the API.
# Used locally (./scripts/verify.sh) and by the GitHub workflow.
set -euo pipefail
cd "$(dirname "$0")/.."

API_PID=""
cleanup() {
  [ -n "$API_PID" ] && kill "$API_PID" 2>/dev/null || true
  docker compose down -v >/dev/null 2>&1 || true
}
trap cleanup EXIT

step() {
  echo
  echo "==> $1"
}

fail() {
  echo "FAIL: $1" >&2
  exit 1
}

# Assert that a curl call returns the expected HTTP status.
# usage: expect_status <expected> <description> <curl args...>
expect_status() {
  local expected=$1 description=$2
  shift 2
  local actual
  actual=$(curl -s -o /dev/null -w '%{http_code}' "$@")
  if [ "$actual" != "$expected" ]; then
    fail "$description: expected HTTP $expected, got $actual"
  fi
  echo "ok: $description ($actual)"
}

step "gofmt"
unformatted=$(gofmt -l .)
if [ -n "$unformatted" ]; then
  echo "$unformatted"
  fail "files need gofmt"
fi

step "go vet"
go vet ./...

step "go build"
go build -o ./bin/api ./cmd/api
go build -o ./bin/migrate ./cmd/migrate

step "go test"
go test ./...

step "migration cycle against a fresh Postgres"
docker compose down -v >/dev/null 2>&1 || true
docker compose up -d db --wait

# The app and the migrate CLI read config from the environment; use the
# committed dev .env but talk to Postgres through the published port.
set -a
source ./.env
set +a
export DB_HOST=localhost
export SERVER_PORT=18080

./bin/migrate up
./bin/migrate down
./bin/migrate up

step "API smoke test"
./bin/api &
API_PID=$!

base="http://localhost:$SERVER_PORT"
for _ in $(seq 1 50); do
  curl -sf "$base/livez" >/dev/null 2>&1 && break
  kill -0 "$API_PID" 2>/dev/null || fail "API process exited during startup"
  sleep 0.2
done
curl -sf "$base/livez" >/dev/null || fail "API did not become healthy"

cookie_jar=$(mktemp)
user='{"name":"Ana","lastName":"Silva","username":"ana","password":"secret123"}'

expect_status 201 "register user" \
  -X POST "$base/v1/users" -d "$user"
expect_status 409 "duplicate username rejected" \
  -X POST "$base/v1/users" -d "$user"
expect_status 401 "wrong password rejected" \
  -X POST "$base/v1/users/login" -d '{"username":"ana","password":"wrong"}'
expect_status 200 "login" \
  -c "$cookie_jar" -X POST "$base/v1/users/login" -d '{"username":"ana","password":"secret123"}'
expect_status 401 "auth required without session" \
  "$base/v1/products"
expect_status 201 "create product" \
  -b "$cookie_jar" -X POST "$base/v1/products" -d '{"name":"Milk","description":"2L"}'
expect_status 200 "list products" \
  -b "$cookie_jar" "$base/v1/products"
expect_status 201 "create group" \
  -b "$cookie_jar" -X POST "$base/v1/groups" -d '{"name":"Home"}'
expect_status 422 "validation error on empty group name" \
  -b "$cookie_jar" -X POST "$base/v1/groups" -d '{"name":""}'

rm -f "$cookie_jar"
echo
echo "All checks passed."
