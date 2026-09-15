#!/usr/bin/env bash
# Checks that everything needed to build, test, and run the project is
# installed. Prints install hints for anything missing. Safe to run anytime:
# it only reads, never installs.
set -uo pipefail

failures=0

ok() {
  echo "ok: $1"
}

missing() {
  echo "MISSING: $1" >&2
  echo "         $2" >&2
  failures=$((failures + 1))
}

# Go: any release >= 1.21 works, because the toolchain line in go.mod makes
# the go command download the exact pinned version automatically.
if command -v go >/dev/null 2>&1; then
  go_minor=$(go env GOVERSION | sed -E 's/^go1\.([0-9]+).*/\1/')
  if [ "${go_minor:-0}" -ge 21 ] 2>/dev/null; then
    ok "go ($(go env GOVERSION))"
  else
    missing "go >= 1.21 (found $(go env GOVERSION))" \
      "upgrade via your package manager or https://go.dev/dl/"
  fi
else
  missing "go" \
    "install from https://go.dev/dl/ (or: brew install go / pacman -S go / apt install golang)"
fi

# Docker engine
if command -v docker >/dev/null 2>&1; then
  if docker info >/dev/null 2>&1; then
    ok "docker ($(docker --version | sed 's/Docker version //;s/,.*//'))"
  else
    missing "a running docker daemon" \
      "docker is installed but the daemon is not reachable; start Docker Desktop or: sudo systemctl start docker"
  fi
else
  missing "docker" \
    "install from https://docs.docker.com/get-docker/"
fi

# Docker Compose v2 (the 'docker compose' subcommand, not legacy docker-compose)
if docker compose version >/dev/null 2>&1; then
  ok "docker compose ($(docker compose version --short 2>/dev/null))"
else
  missing "docker compose v2" \
    "ships with Docker Desktop; on Linux: https://docs.docker.com/compose/install/linux/"
fi

# curl: used by the smoke test in verify.sh
if command -v curl >/dev/null 2>&1; then
  ok "curl"
else
  missing "curl" \
    "install via your package manager (brew install curl / pacman -S curl / apt install curl)"
fi

if [ "$failures" -gt 0 ]; then
  echo >&2
  echo "$failures prerequisite(s) missing." >&2
  exit 1
fi

echo
echo "All prerequisites present."
