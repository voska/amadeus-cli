# amadeus — Amadeus Travel API CLI

[![CI](https://github.com/voska/amadeus-cli/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/voska/amadeus-cli/actions/workflows/ci.yml)
[![Go](https://img.shields.io/github/go-mod/go-version/voska/amadeus-cli)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

CLI for humans and AI agents. Data goes to stdout (parseable), hints/progress to stderr.

```bash
$ amadeus flights search --from JFK --to CDG --date 2026-04-15 --json
[
  {"id": "1", "price": {"total": "54.74", "currency": "EUR"}, "itineraries": [...]},
  {"id": "2", "price": {"total": "121.33", "currency": "EUR"}, "itineraries": [...]}
]

$ amadeus airports search london
name       iataCode  subType
─────────  ────────  ───────
HEATHROW   LHR       AIRPORT
GATWICK    LGW       AIRPORT
STANSTED   STN       AIRPORT
```

Run `amadeus --help` for the full command tree, or `amadeus schema` for machine-readable introspection.

## Install

**Homebrew** (macOS / Linux):

```bash
brew install voska/tap/amadeus
```

**Go**:

```bash
go install github.com/voska/amadeus-cli/cmd/amadeus@latest
```

**Binary**: download from [Releases](https://github.com/voska/amadeus-cli/releases).

## Getting Credentials

1. Sign up at [developers.amadeus.com](https://developers.amadeus.com/)
2. Create an app in the dashboard to get your **API Key** and **API Secret**
3. The test environment is free and rate-limited with synthetic data

## Quick Start

```bash
export AMADEUS_API_KEY=your_key
export AMADEUS_API_SECRET=your_secret

# Authenticate (test environment)
amadeus auth login --test

# Verify
amadeus auth status --test --json

# Search flights
amadeus flights search --from JFK --to CDG --date 2026-04-15 --test

# Search hotels in Paris
amadeus hotels search --city PAR --test

# Get hotel offers
amadeus hotels offers --hotel-id HLPAR123 --checkin 2026-04-15 --checkout 2026-04-18 --test

# Look up airports and airlines
amadeus airports search london --test
amadeus airlines lookup BA --test
```

## Output Modes

| Flag | Description |
|------|-------------|
| (default) | Colored tables, summaries on stderr |
| `--json` / `-j` | Structured JSON to stdout |
| `--plain` / `-p` | Tab-separated values, no color |
| `--results-only` | Strip response metadata, return data array |
| `--select f1,f2` | Project output to specific fields (dot-path) |

Auto-JSON: when stdout is not a TTY and `AMADEUS_AUTO_JSON=1`, defaults to JSON output.

## Commands

| Command | Description |
|---------|-------------|
| `auth login\|status\|logout` | OAuth2 authentication |
| `flights search` | Search flight offers |
| `flights price` | Confirm pricing for an offer |
| `hotels search` | Search hotels by city or geocode |
| `hotels offers` | Get offers for a specific hotel |
| `hotels book` | Book a hotel offer |
| `airports search <keyword>` | Airport autocomplete search |
| `airlines lookup <code>` | Airline lookup by IATA code |
| `schema [command]` | CLI command tree as JSON |
| `exit-codes` | Exit code reference |

All commands support `--dry-run`, `--no-input`, and `--test`. Run `amadeus exit-codes` for the full exit code reference.

## Environments

- **Test** (default): `test.api.amadeus.com` — free, rate-limited, synthetic data
- **Production**: `api.amadeus.com` — requires a paid plan

Use `--test` or set `environment: "production"` in `~/.config/amadeus/config.json`.

## Development

```bash
make build    # Build to bin/amadeus
make test     # Run tests
make lint     # Run linter
make vet      # Run go vet
make fmt      # Format code
make generate # Regenerate API types from OpenAPI specs
```

## License

MIT
