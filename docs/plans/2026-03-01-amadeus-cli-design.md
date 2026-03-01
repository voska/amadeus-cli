# amadeus-cli Design

## Overview

A CLI for the Amadeus travel APIs, designed primarily for AI agent consumption. Follows the same architecture and patterns as [qbo-cli](https://github.com/voska/qbo-cli).

No dedicated Amadeus CLI exists today — only SDKs (Node.js, Python, Java). This fills that gap.

## Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Name | `amadeus` | Direct, obvious, no conflict for Go binary |
| Language | Go | Single binary, fast startup, matches qbo-cli |
| CLI framework | Kong | Declarative struct-based, matches qbo-cli |
| API client | oapi-codegen | Generate typed client from Amadeus OpenAPI specs |
| Auth | Client credentials OAuth2 | Env vars + keyring token cache |
| Distribution | GoReleaser + Homebrew | Same as qbo-cli |
| v1 scope | Flights + Hotels | ~10 endpoints, expand later |

## Architecture

```
cmd/amadeus/main.go          # Thin entrypoint, wires Kong
internal/
  cmd/                        # Kong command structs
    root.go                   # CLI struct, Globals
    auth.go                   # auth login, auth status, auth logout
    flights.go                # flights search, flights price, flights seatmap
    airports.go               # airports search (autocomplete)
    airlines.go               # airlines lookup
    hotels.go                 # hotels search, hotels offers, hotels book
    schema.go                 # Agent introspection (schema --json)
    exit_codes.go             # Structured exit codes
  api/                        # Generated + hand-written HTTP client
    client.go                 # Base HTTP client, token management
    generated/                # oapi-codegen output (types + client)
      flights.go
      hotels.go
  auth/                       # OAuth2 client credentials
    oauth.go                  # Token endpoint, auto-refresh
    token.go                  # Keyring storage
    config.go                 # ~/.config/amadeus/ credentials
  output/                     # Output formatting (reuse qbo-cli patterns)
    mode.go                   # JSON / Plain / Human
    write.go
    human.go
    plain.go
    json.go
    select.go
  config/
    config.go                 # ~/.config/amadeus/config.json
  errfmt/
    errors.go                 # Exit codes
```

## v1 Commands

```
amadeus auth login            # Store API key + secret, fetch initial token
amadeus auth status           # Show token expiry, environment (test/prod)
amadeus auth logout           # Clear keyring

amadeus flights search        # --from JFK --to CDG --date 2026-04-01 --adults 1
amadeus flights price         # --offer <offer-id>
amadeus flights seatmap       # --offer <offer-id>

amadeus hotels search         # --city PAR (or --lat/--lng)
amadeus hotels offers         # --hotel-id XXXX --checkin --checkout
amadeus hotels book           # --offer-id XXXX --guest-name "..." --guest-email "..."

amadeus airports search       # --keyword "london" (autocomplete)
amadeus airlines lookup       # --code BA

amadeus schema                # Full CLI introspection for agents
amadeus exit-codes            # Structured exit code reference
```

## Agent-Friendly Features

- `--json` / `-j` on all commands
- `--plain` / `-p` for TSV
- `--results-only` to unwrap response envelopes
- `--select field1,field2` for field projection
- `--dry-run` on mutations (book)
- `--no-input` for non-interactive agent use
- `AMADEUS_AUTO_JSON=1` auto-enables JSON when piped
- Structured exit codes (0=success, 3=empty, 4=auth, 5=not found, 7=rate limited)
- `schema --json` for full CLI tree introspection
- Data to stdout, hints/errors to stderr

## Auth Flow

1. User sets `AMADEUS_API_KEY` + `AMADEUS_API_SECRET` (env vars or `~/.config/amadeus/config.json`)
2. `amadeus auth login` posts to `/v1/security/oauth2/token` with client credentials
3. Bearer token cached in OS keyring (macOS Keychain, Linux Secret Service, Windows Credential Manager)
4. Auto-refreshed on expiry (Amadeus tokens last 30 minutes)
5. `--test` flag switches between `test.api.amadeus.com` and `api.amadeus.com`

## API Client Strategy

- Amadeus publishes OpenAPI specs at [amadeus4dev/amadeus-open-api-specification](https://github.com/amadeus4dev/amadeus-open-api-specification)
- Use [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) to generate Go types and client code
- Specs may be Swagger 2.0 (filenames say "swagger") — convert to OpenAPI 3.0 if needed
- Generated code lives in `internal/api/generated/`
- Hand-written `client.go` wraps generated client with auth, error handling, retries

## Distribution

- GoReleaser for cross-platform builds (Linux, macOS, Windows; amd64, arm64)
- Homebrew tap
- Binary downloads
- Claude Code skill: `npx skills add -g voska/amadeus-cli`

## Dependencies

- `github.com/alecthomas/kong` — CLI framework
- `github.com/99designs/keyring` — Secure token storage
- `github.com/muesli/termenv` — Terminal colors
- `github.com/oapi-codegen/oapi-codegen` — OpenAPI code generation (dev dependency)
- `golang.org/x/oauth2` — OAuth2 (or hand-rolled, since client credentials is trivial)

## Future Scope (post-v1)

- Cars & transfers
- Destination experiences / points of interest
- Market insights (flight traffic, busiest periods)
- Trip planning (itinerary management)
- Booking management (retrieve, cancel)
