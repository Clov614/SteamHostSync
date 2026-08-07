# CLAUDE.md — Project guide for SteamHostSync

## Project Overview

SteamHostSync is a **Go 1.22+** CLI that generates hosts files for accelerating Steam, GitHub, Docker, GOG, and Ubisoft. It resolves each platform's domains via multiple DoH resolvers, probes candidate IPs over TCP (default port 443), picks the fastest reachable IP, and writes per-platform `Hosts_<name>` files plus a combined `Hosts` and a README.

## Development Workflow (branch convention — follow this for every change)

- **Open a dedicated branch from clean `main`** for each new feature/fix/docs change. Do NOT build features on top of an unrelated in-progress branch.
- Recommended branch naming (mirrors conventional commits):
  - `feat/<feature>` — new functionality
  - `fix/<fix>` — bug fixes
  - `docs/<doc>` — documentation only
- Complete the work on that branch: TDD → implement → verify (`go vet` + `go test`) → commit (conventional message) → push → open a PR back to `main`.
- Keep `main` releasable. The CI auto-commit (`[skip ci] auto update hosts`) only runs on `main` flow; it regenerates the `Hosts*` and `README.md` artifacts — do not commit those generated artifacts manually.

## Architecture & Key Files

- `main.go` — entrypoint / CLI flag parsing.
- `internal/app/app.go` — orchestrates the resolve → probe → render → write pipeline; degrades gracefully when domains or whole platforms fail.
- `internal/config/` — config parsing & validation; `default_config.yaml` is the embedded default (single source for a missing `config.yaml`). Platforms are purely config-driven: each `platforms` entry produces one `Hosts_<name>` file.
- `internal/resolve/` — DoH resolution with multi-resolver fallback + blacklist.
- `internal/probe/` — TCP probing to select the best IP.
- `internal/render/` — hosts content generation (pure functions).
- `internal/fileio/` — atomic writes (temp file + rename).

## Testing & Quality gates

- Coverage must stay ≥ 80% (CI enforces this).
- Run locally: `go vet ./... && go test -cover ./internal/...`
- The CI runs `go test -race ./internal/...` on Linux; `-race` needs CGO and may not build with the local Windows toolchain — run plain `go test` locally.

## Platform change example

To add/expose a host category (e.g. a Linux-specific Steam repo), add a new `platforms` entry to **both** `config.yaml` and `internal/config/default_config.yaml`, add/verify the parsed domains, and update the docs (`README_TEMP.md` feeds the generated `README.md`; `README.en.md` / `README.zh-CN.md` are manual).