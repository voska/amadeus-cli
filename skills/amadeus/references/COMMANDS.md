# amadeus-cli Command Reference

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--json` | `-j` | Output JSON to stdout (alias: `--machine`) |
| `--plain` | `-p` | Output TSV, no color |
| `--select` | `-s` | Project output to fields using dot-path (alias: `--fields`) |
| `--results-only` | | Strip response metadata, return data array only |
| `--test` | | Use test environment (test.api.amadeus.com) |
| `--dry-run` | | Show what would happen without making API calls |
| `--no-input` | | Never prompt for input |
| `--verbose` | `-v` | Verbose output to stderr |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `AMADEUS_API_KEY` | API key (overrides config file) |
| `AMADEUS_API_SECRET` | API secret (overrides config file) |
| `AMADEUS_AUTO_JSON` | Set to `1` to auto-detect non-TTY and output JSON |
| `AMADEUS_CONFIG_DIR` | Override config directory (default: `~/.config/amadeus`) |

## Commands

### `amadeus auth login`

Authenticate with Amadeus API. Requires `AMADEUS_API_KEY` and `AMADEUS_API_SECRET` to be set.

### `amadeus auth status`

Show current authentication status (environment, expiry, token validity).

### `amadeus auth logout`

Remove stored credentials from keyring.

### `amadeus flights search`

Search for flight offers.

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--from` | yes | | Origin IATA code (e.g., JFK) |
| `--to` | yes | | Destination IATA code (e.g., CDG) |
| `--date` | yes | | Departure date (YYYY-MM-DD) |
| `--return` | | | Return date for round-trip (YYYY-MM-DD) |
| `--adults` | | 1 | Number of adults (1-9) |
| `--children` | | 0 | Number of children (0-9) |
| `--class` | | | Travel class: ECONOMY, PREMIUM_ECONOMY, BUSINESS, FIRST |
| `--nonstop` | | false | Direct flights only |
| `--currency` | | | Currency code (e.g., USD, EUR) |
| `--max-price` | | | Maximum price (no decimals) |
| `--max` | | 10 | Maximum number of results (1-250) |

API: `GET /v2/shopping/flight-offers`

### `amadeus flights price`

Confirm pricing for a flight offer. *(Not yet implemented)*

| Flag | Required | Description |
|------|----------|-------------|
| `--offer-id` | yes | Flight offer ID from search results |

### `amadeus flights seatmap`

View seat map for a flight offer. *(Not yet implemented)*

| Flag | Required | Description |
|------|----------|-------------|
| `--offer-id` | yes | Flight offer ID from search results |

### `amadeus hotels search`

Search hotels by city or geocode location.

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--city` | xor:location | | City IATA code (e.g., PAR) |
| `--lat` | xor:location | | Latitude |
| `--lng` | xor:location | | Longitude |
| `--radius` | | 5 | Search radius in km |
| `--ratings` | | | Hotel ratings to filter (1-5, comma-separated) |

API: `GET /v1/reference-data/locations/hotels/by-city` or `GET /v1/reference-data/locations/hotels/by-geocode`

### `amadeus hotels offers`

Get offers for a specific hotel.

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--hotel-id` | yes | | Amadeus hotel ID |
| `--checkin` | yes | | Check-in date (YYYY-MM-DD) |
| `--checkout` | yes | | Check-out date (YYYY-MM-DD) |
| `--adults` | | 1 | Number of adults |
| `--rooms` | | 1 | Number of rooms |
| `--currency` | | | Currency code |

API: `GET /v3/shopping/hotel-offers`

### `amadeus hotels book`

Book a hotel offer.

| Flag | Required | Description |
|------|----------|-------------|
| `--offer-id` | yes | Hotel offer ID from search results |
| `--guest-name` | yes | Guest full name |
| `--guest-email` | yes | Guest email address |

API: `POST /v2/booking/hotel-orders`

### `amadeus airports search <keyword>`

Search airports by keyword (autocomplete).

| Argument | Required | Description |
|----------|----------|-------------|
| `keyword` | yes | Search keyword (e.g., "london", "JFK") |

API: `GET /v1/reference-data/locations?subType=AIRPORT`

### `amadeus airlines lookup <code>`

Look up airline by IATA code.

| Argument | Required | Description |
|----------|----------|-------------|
| `code` | yes | IATA airline code (e.g., BA, AA) |

API: `GET /v1/reference-data/airlines`

### `amadeus schema [command]`

Show CLI schema as JSON for agent introspection. Optionally filter to a specific command.

### `amadeus exit-codes`

Show exit code reference table.
