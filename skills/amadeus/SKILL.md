# amadeus-cli — Amadeus Travel API CLI

Agent-friendly CLI for searching flights, hotels, airports, and airlines via the Amadeus travel APIs.

## Setup

### Getting API Credentials

1. Sign up for a free account at https://developers.amadeus.com/
2. Go to **My Self-Service Workspace** → **Create New App**
3. Copy the **API Key** and **API Secret** from the app dashboard

The test environment is free, rate-limited, and uses synthetic data. No credit card required.

### Configure and Authenticate

```bash
export AMADEUS_API_KEY=your_key
export AMADEUS_API_SECRET=your_secret
amadeus auth login
```

If the user doesn't have credentials yet, direct them to https://developers.amadeus.com/ to sign up and create an app.

## Agent Integration

- Use `--json` (or `-j`) for structured output
- Use `--results-only` to strip metadata and get the `data` array directly
- Use `--select field1,field2` to project specific fields (supports dot-paths like `price.total`)
- Use `--dry-run` to preview API calls without making them
- Set `AMADEUS_AUTO_JSON=1` to auto-detect piped output and switch to JSON
- Check exit codes programmatically (see `amadeus exit-codes --json`)

## Introspection

```bash
amadeus schema --json          # Full CLI tree with all commands, flags, and args
amadeus schema flights --json  # Schema for a specific command
amadeus exit-codes --json      # Exit code reference
```

## Common Workflows

### Search flights

```bash
# One-way JFK to CDG
amadeus flights search --from JFK --to CDG --date 2026-04-15 --json

# Round-trip, business class, max 5 results
amadeus flights search --from JFK --to CDG --date 2026-04-15 --return 2026-04-22 --class BUSINESS --max 5 --json

# Direct flights only, with price filter
amadeus flights search --from LAX --to NRT --date 2026-05-01 --nonstop --max-price 2000 --json

# Get just prices
amadeus flights search --from JFK --to CDG --date 2026-04-15 --results-only --select id,price.total,price.currency --json
```

### Search and book hotels

```bash
# Find hotels in Paris
amadeus hotels search --city PAR --json

# Find hotels near coordinates
amadeus hotels search --lat 48.8566 --lng 2.3522 --radius 10 --json

# Get offers for a specific hotel
amadeus hotels offers --hotel-id HLPAR123 --checkin 2026-04-15 --checkout 2026-04-18 --json

# Book a hotel
amadeus hotels book --offer-id ABC123 --guest-name "John Doe" --guest-email john@example.com --json
```

### Look up airports and airlines

```bash
# Search airports by keyword
amadeus airports search london --json

# Look up airline by IATA code
amadeus airlines lookup BA --json
```

## Environments

- **Test** (default): `test.api.amadeus.com` — free, rate-limited, synthetic data
- **Production**: `api.amadeus.com` — requires paid plan

Use `--test` flag or set `environment: "test"` in `~/.config/amadeus/config.json`.

## Exit Codes

| Code | Name | Meaning |
|------|------|---------|
| 0 | success | Command completed |
| 1 | error | General error |
| 2 | usage | Invalid arguments |
| 3 | empty | No results found |
| 4 | auth_required | Not authenticated |
| 5 | not_found | HTTP 404 |
| 6 | forbidden | HTTP 403 |
| 7 | rate_limited | HTTP 429 |
| 8 | retryable | HTTP 5xx |
| 10 | config_error | Config problem |
