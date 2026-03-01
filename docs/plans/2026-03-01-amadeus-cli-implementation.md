# amadeus-cli Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build an agent-friendly CLI for Amadeus travel APIs (flights + hotels), mirroring qbo-cli architecture.

**Architecture:** Go binary using Kong for CLI parsing, oapi-codegen for typed API models (Swagger 2.0 specs converted to OpenAPI 3.0), keyring for token storage, structured output (JSON/plain/human) with agent-friendly features (exit codes, schema introspection, field projection).

**Tech Stack:** Go 1.25, Kong, oapi-codegen, 99designs/keyring, termenv, GoReleaser

**Reference:** Mirror patterns from `/Users/matt/Developer/qbo-cli` throughout.

---

### Task 1: Project scaffolding

**Files:**
- Create: `go.mod`
- Create: `cmd/amadeus/main.go`
- Create: `Makefile`
- Create: `.gitignore`

**Step 1: Initialize Go module**

Run: `cd /Users/matt/Developer/amadeus-cli && go mod init github.com/voska/amadeus-cli`

**Step 2: Create .gitignore**

```
bin/
dist/
*.exe
.DS_Store
```

**Step 3: Create minimal main.go**

```go
// cmd/amadeus/main.go
package main

import (
	"fmt"
	"os"
)

var version = "dev"

func main() {
	fmt.Fprintf(os.Stderr, "amadeus %s\n", version)
	os.Exit(0)
}
```

**Step 4: Create Makefile**

```makefile
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X main.version=$(VERSION)
BIN     := bin/amadeus

.PHONY: build test lint clean install generate

build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/amadeus

test:
	go test -race ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ dist/

install: build
	cp $(BIN) $(GOPATH)/bin/amadeus 2>/dev/null || cp $(BIN) ~/go/bin/amadeus

generate:
	go generate ./...
```

**Step 5: Verify build**

Run: `make build && ./bin/amadeus`
Expected: `amadeus dev`

**Step 6: Commit**

```bash
git add -A && git commit -m "feat: project scaffolding with Go module, main entrypoint, Makefile"
```

---

### Task 2: Error types and exit codes

**Files:**
- Create: `internal/errfmt/errors.go`
- Create: `internal/cmd/exit_codes.go`

**Step 1: Create error types**

Reference: `/Users/matt/Developer/qbo-cli/internal/errfmt/errors.go`

```go
// internal/errfmt/errors.go
package errfmt

import "fmt"

const (
	ExitOK        = 0
	ExitError     = 1
	ExitUsage     = 2
	ExitEmpty     = 3
	ExitAuth      = 4
	ExitNotFound  = 5
	ExitForbidden = 6
	ExitRateLimit = 7
	ExitRetryable = 8
	ExitConfig    = 10
)

type Error struct {
	Code    int
	Message string
	Detail  string
}

func (e *Error) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Detail)
	}
	return e.Message
}

func New(code int, msg string) *Error {
	return &Error{Code: code, Message: msg}
}

func Wrap(code int, msg string, err error) *Error {
	return &Error{Code: code, Message: msg, Detail: err.Error()}
}

func Auth(msg string) *Error      { return New(ExitAuth, msg) }
func NotFound(msg string) *Error   { return New(ExitNotFound, msg) }
func Usage(msg string) *Error      { return New(ExitUsage, msg) }
func Empty() *Error                { return New(ExitEmpty, "no results") }
func Config(msg string) *Error     { return New(ExitConfig, msg) }
func RateLimit() *Error            { return New(ExitRateLimit, "rate limited") }
func Forbidden(msg string) *Error  { return New(ExitForbidden, msg) }
```

**Step 2: Create exit codes command**

Reference: `/Users/matt/Developer/qbo-cli/internal/cmd/exit_codes.go`

```go
// internal/cmd/exit_codes.go
package cmd

import "github.com/voska/amadeus-cli/internal/errfmt"

type ExitCodeEntry struct {
	Code        int    `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ExitCodesCmd struct{}

func ExitCodeTable() []ExitCodeEntry {
	return []ExitCodeEntry{
		{errfmt.ExitOK, "success", "Command completed successfully"},
		{errfmt.ExitError, "error", "General error"},
		{errfmt.ExitUsage, "usage", "Invalid arguments or usage"},
		{errfmt.ExitEmpty, "empty", "No results found"},
		{errfmt.ExitAuth, "auth_required", "Authentication required or token expired"},
		{errfmt.ExitNotFound, "not_found", "Resource not found (HTTP 404)"},
		{errfmt.ExitForbidden, "forbidden", "Permission denied (HTTP 403)"},
		{errfmt.ExitRateLimit, "rate_limited", "Rate limited (HTTP 429)"},
		{errfmt.ExitRetryable, "retryable", "Transient error, retry may succeed (HTTP 5xx)"},
		{errfmt.ExitConfig, "config_error", "Configuration error"},
	}
}
```

**Step 3: Verify compilation**

Run: `go build ./...`
Expected: No errors

**Step 4: Commit**

```bash
git add internal/errfmt/ internal/cmd/exit_codes.go && git commit -m "feat: structured error types and exit codes"
```

---

### Task 3: Output formatting layer

**Files:**
- Create: `internal/output/mode.go`
- Create: `internal/output/write.go`
- Create: `internal/output/json.go`
- Create: `internal/output/human.go`
- Create: `internal/output/plain.go`
- Create: `internal/output/select.go`

**Step 1: Create output mode and options**

Reference: `/Users/matt/Developer/qbo-cli/internal/output/mode.go`

```go
// internal/output/mode.go
package output

import (
	"context"
	"os"
)

type Mode int

const (
	ModeHuman Mode = iota
	ModeJSON
	ModePlain
)

type Options struct {
	Mode        Mode
	ResultsOnly bool
	Select      []string
	Pretty      bool
}

type ctxKey struct{}

func WithOptions(ctx context.Context, opts Options) context.Context {
	return context.WithValue(ctx, ctxKey{}, opts)
}

func GetOptions(ctx context.Context) Options {
	if opts, ok := ctx.Value(ctxKey{}).(Options); ok {
		return opts
	}
	return Options{Mode: ModeHuman, Pretty: true}
}

func isTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
```

**Step 2: Create JSON writer**

Reference: `/Users/matt/Developer/qbo-cli/internal/output/json.go`

```go
// internal/output/json.go
package output

import (
	"encoding/json"
	"os"
)

func WriteJSON(data any, pretty bool) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	if pretty {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(data)
}
```

**Step 3: Create human-readable writer with colors**

Reference: `/Users/matt/Developer/qbo-cli/internal/output/human.go`

Implement colored table output with termenv, plus Hint/Success/Warn/ErrorMsg functions that write to stderr.

**Step 4: Create plain writer**

Reference: `/Users/matt/Developer/qbo-cli/internal/output/plain.go`

Tab-separated output, no color, headers on first line.

**Step 5: Create field projection**

Reference: `/Users/matt/Developer/qbo-cli/internal/output/select.go`

Implement `ProjectFields(data any, fields []string) any` with dot-path support and `StripMetadata(data any) any` for unwrapping Amadeus response envelopes (strip `meta` and `dictionaries` keys, return `data` array).

**Step 6: Create main Write dispatcher**

Reference: `/Users/matt/Developer/qbo-cli/internal/output/write.go`

```go
// internal/output/write.go
package output

import "context"

func Write(ctx context.Context, data any) error {
	opts := GetOptions(ctx)
	if opts.ResultsOnly {
		data = StripMetadata(data)
	}
	if len(opts.Select) > 0 {
		data = ProjectFields(data, opts.Select)
	}
	switch opts.Mode {
	case ModeJSON:
		return WriteJSON(data, opts.Pretty)
	case ModePlain:
		headers, rows := toTable(data)
		return WritePlain(headers, rows)
	default:
		headers, rows := toTable(data)
		return WriteTable(headers, rows)
	}
}
```

**Step 7: Add termenv dependency**

Run: `go get github.com/muesli/termenv@v0.16.0`

**Step 8: Verify compilation**

Run: `go build ./...`

**Step 9: Commit**

```bash
git add internal/output/ go.mod go.sum && git commit -m "feat: output formatting layer (JSON, plain, human with colors)"
```

---

### Task 4: Config management

**Files:**
- Create: `internal/config/config.go`

**Step 1: Create config types and file management**

Reference: `/Users/matt/Developer/qbo-cli/internal/config/config.go`

```go
// internal/config/config.go
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey      string `json:"api_key,omitempty"`
	APISecret   string `json:"api_secret,omitempty"`
	Environment string `json:"environment,omitempty"` // "test" or "production"
}

func Dir() (string, error) {
	if d := os.Getenv("AMADEUS_CONFIG_DIR"); d != "" {
		return d, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "amadeus"), nil
}

func Load() (*Config, error) {
	dir, err := Dir()
	if err != nil {
		return &Config{}, err
	}
	data, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if os.IsNotExist(err) {
		return &Config{}, nil
	}
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "config.json"), data, 0o600)
}

// ResolveAPIKey returns API key from env var or config file.
func (c *Config) ResolveAPIKey() string {
	if v := os.Getenv("AMADEUS_API_KEY"); v != "" {
		return v
	}
	return c.APIKey
}

// ResolveAPISecret returns API secret from env var or config file.
func (c *Config) ResolveAPISecret() string {
	if v := os.Getenv("AMADEUS_API_SECRET"); v != "" {
		return v
	}
	return c.APISecret
}

// ResolveEnvironment returns "test" or "production".
func (c *Config) ResolveEnvironment(flagTest bool) string {
	if flagTest {
		return "test"
	}
	if c.Environment != "" {
		return c.Environment
	}
	return "test" // default to test for safety
}

// BaseURL returns the API base URL for the environment.
func BaseURL(env string) string {
	if env == "production" {
		return "https://api.amadeus.com"
	}
	return "https://test.api.amadeus.com"
}
```

**Step 2: Verify compilation**

Run: `go build ./...`

**Step 3: Commit**

```bash
git add internal/config/ && git commit -m "feat: config management with env var + file resolution"
```

---

### Task 5: Auth (OAuth2 client credentials + keyring)

**Files:**
- Create: `internal/auth/oauth.go`
- Create: `internal/auth/token.go`

**Step 1: Create OAuth2 client credentials flow**

Amadeus uses a simple client_credentials grant. No browser redirect needed — just POST key+secret, get bearer token.

```go
// internal/auth/oauth.go
package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/voska/amadeus-cli/internal/config"
)

type Token struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int       `json:"expires_in"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// FetchToken exchanges API key + secret for a bearer token.
func FetchToken(apiKey, apiSecret, env string) (*Token, error) {
	base := config.BaseURL(env)
	endpoint := base + "/v1/security/oauth2/token"

	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {apiKey},
		"client_secret": {apiSecret},
	}

	resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request returned %d: %s", resp.StatusCode, string(body))
	}

	var tok Token
	if err := json.Unmarshal(body, &tok); err != nil {
		return nil, fmt.Errorf("parsing token response: %w", err)
	}
	tok.ExpiresAt = time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
	return &tok, nil
}
```

**Step 2: Create keyring token storage**

Reference: `/Users/matt/Developer/qbo-cli/internal/auth/token.go`

```go
// internal/auth/token.go
package auth

import (
	"encoding/json"

	"github.com/99designs/keyring"
)

const serviceName = "amadeus-cli"

func openKeyring() (keyring.Keyring, error) {
	dir, err := configDir()
	if err != nil {
		return nil, err
	}
	return keyring.Open(keyring.Config{
		ServiceName:      serviceName,
		FileDir:          dir + "/tokens",
		FilePasswordFunc: func(string) (string, error) { return "", nil },
	})
}

func StoreToken(env string, token *Token) error {
	kr, err := openKeyring()
	if err != nil {
		return err
	}
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return kr.Set(keyring.Item{
		Key:  env,
		Data: data,
	})
}

func LoadToken(env string) (*Token, error) {
	kr, err := openKeyring()
	if err != nil {
		return nil, err
	}
	item, err := kr.Get(env)
	if err != nil {
		return nil, err
	}
	var tok Token
	if err := json.Unmarshal(item.Data, &tok); err != nil {
		return nil, err
	}
	return &tok, nil
}

func DeleteToken(env string) error {
	kr, err := openKeyring()
	if err != nil {
		return err
	}
	return kr.Remove(env)
}
```

Use the config package's Dir() for `configDir()`.

**Step 3: Add keyring dependency**

Run: `go get github.com/99designs/keyring@v1.2.2`

**Step 4: Verify compilation**

Run: `go build ./...`

**Step 5: Commit**

```bash
git add internal/auth/ go.mod go.sum && git commit -m "feat: OAuth2 client credentials auth with keyring token storage"
```

---

### Task 6: HTTP API client

**Files:**
- Create: `internal/api/client.go`

**Step 1: Create the API client**

Reference: `/Users/matt/Developer/qbo-cli/internal/api/client.go`

```go
// internal/api/client.go
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/voska/amadeus-cli/internal/auth"
	"github.com/voska/amadeus-cli/internal/errfmt"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	token      *auth.Token
}

func NewClient(token *auth.Token, env string) *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    baseURL(env),
		token:      token,
	}
}

func baseURL(env string) string {
	if env == "production" {
		return "https://api.amadeus.com"
	}
	return "https://test.api.amadeus.com"
}

// Get makes an authenticated GET request. path should include version prefix (e.g., "/v2/shopping/flight-offers").
func (c *Client) Get(path string, params url.Values) (map[string]any, error) {
	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

// Post makes an authenticated POST request with JSON body.
func (c *Client) Post(path string, body any) (map[string]any, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.baseURL+path, strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/vnd.amadeus+json")
	return c.do(req)
}

func (c *Client) do(req *http.Request) (map[string]any, error) {
	req.Header.Set("Authorization", "Bearer "+c.token.AccessToken)
	req.Header.Set("Accept", "application/vnd.amadeus+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return result, nil
	case http.StatusUnauthorized:
		return nil, errfmt.Auth("authentication failed — run 'amadeus auth login'")
	case http.StatusForbidden:
		return nil, errfmt.Forbidden(extractError(result))
	case http.StatusNotFound:
		return nil, errfmt.NotFound(extractError(result))
	case http.StatusTooManyRequests:
		return nil, errfmt.RateLimit()
	default:
		if resp.StatusCode >= 500 {
			return nil, errfmt.New(errfmt.ExitRetryable, extractError(result))
		}
		return nil, errfmt.New(errfmt.ExitError, extractError(result))
	}
}

func extractError(result map[string]any) string {
	// Amadeus error format: {"errors": [{"detail": "..."}]}
	if errors, ok := result["errors"].([]any); ok && len(errors) > 0 {
		if first, ok := errors[0].(map[string]any); ok {
			if detail, ok := first["detail"].(string); ok {
				return detail
			}
		}
	}
	return "unknown error"
}
```

**Step 2: Verify compilation**

Run: `go build ./...`

**Step 3: Commit**

```bash
git add internal/api/ && git commit -m "feat: HTTP API client with auth and structured error mapping"
```

---

### Task 7: Download and convert Amadeus OpenAPI specs, generate types

**Files:**
- Create: `specs/` directory with downloaded and converted specs
- Create: `internal/api/generated/` with generated Go types
- Create: `oapi-codegen.yaml` config

**Step 1: Download Amadeus specs**

```bash
mkdir -p specs/swagger specs/openapi
# Flight search
curl -sL https://raw.githubusercontent.com/amadeus4dev/amadeus-open-api-specification/main/spec/yaml/FlightOffersSearch_v2_swagger_specification.yaml -o specs/swagger/flights_search.yaml
# Hotel list (search by city)
curl -sL https://raw.githubusercontent.com/amadeus4dev/amadeus-open-api-specification/main/spec/yaml/HotelList_v1_swagger_specification.yaml -o specs/swagger/hotel_list.yaml
# Hotel search (offers)
curl -sL https://raw.githubusercontent.com/amadeus4dev/amadeus-open-api-specification/main/spec/yaml/HotelSearch_v3_swagger_specification.yaml -o specs/swagger/hotel_search.yaml
# Auth (already OpenAPI 3.0)
curl -sL https://raw.githubusercontent.com/amadeus4dev/amadeus-open-api-specification/main/spec/yaml/Authorizaton_v1_swagger_specification.yaml -o specs/openapi/auth.yaml
```

**Step 2: Convert Swagger 2.0 to OpenAPI 3.0**

```bash
npx swagger2openapi specs/swagger/flights_search.yaml -o specs/openapi/flights_search.yaml
npx swagger2openapi specs/swagger/hotel_list.yaml -o specs/openapi/hotel_list.yaml
npx swagger2openapi specs/swagger/hotel_search.yaml -o specs/openapi/hotel_search.yaml
```

**Step 3: Create oapi-codegen config**

```yaml
# oapi-codegen.yaml
generate:
  - config:
      package: generated
      output: internal/api/generated/flights.go
      generate:
        models: true
        client: false
        embedded-spec: false
    spec: specs/openapi/flights_search.yaml
  - config:
      package: generated
      output: internal/api/generated/hotels.go
      generate:
        models: true
        client: false
        embedded-spec: false
    spec: specs/openapi/hotel_search.yaml
  - config:
      package: generated
      output: internal/api/generated/hotel_list.go
      generate:
        models: true
        client: false
        embedded-spec: false
    spec: specs/openapi/hotel_list.yaml
```

**Step 4: Install oapi-codegen and generate types**

```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

Generate types per spec. If oapi-codegen has trouble with the converted specs, fall back to hand-writing the key types (FlightOffer, HotelOffer, etc.) — the structs are documented in the design doc research.

**Step 5: Verify generated code compiles**

Run: `go build ./...`

**Step 6: Commit**

```bash
git add specs/ internal/api/generated/ oapi-codegen.yaml && git commit -m "feat: download Amadeus OpenAPI specs and generate Go types"
```

---

### Task 8: Kong CLI wiring and Globals

**Files:**
- Create: `internal/cmd/root.go`
- Modify: `cmd/amadeus/main.go`

**Step 1: Create CLI struct and Globals**

Reference: `/Users/matt/Developer/qbo-cli/internal/cmd/root.go`

```go
// internal/cmd/root.go
package cmd

import (
	"context"
	"os"

	"github.com/voska/amadeus-cli/internal/api"
	"github.com/voska/amadeus-cli/internal/auth"
	"github.com/voska/amadeus-cli/internal/config"
	"github.com/voska/amadeus-cli/internal/errfmt"
	"github.com/voska/amadeus-cli/internal/output"
)

type CLI struct {
	// Global flags
	JSON        bool     `short:"j" name:"json" help:"Output as JSON" aliases:"machine"`
	Plain       bool     `short:"p" name:"plain" help:"Output as tab-separated values"`
	ResultsOnly bool     `name:"results-only" help:"Strip response metadata, return data array only"`
	Select      []string `name:"select" short:"s" help:"Select specific fields (dot-path)" aliases:"fields"`
	Test        bool     `name:"test" help:"Use test environment (test.api.amadeus.com)"`
	DryRun      bool     `name:"dry-run" help:"Show what would be done without making API calls"`
	NoInput     bool     `name:"no-input" help:"Never prompt for input"`
	Verbose     bool     `short:"v" name:"verbose" help:"Verbose output"`

	// Commands
	Auth      AuthCmd      `cmd:"" help:"Manage authentication"`
	Flights   FlightsCmd   `cmd:"" help:"Search and price flights"`
	Hotels    HotelsCmd    `cmd:"" help:"Search and book hotels"`
	Airports  AirportsCmd  `cmd:"" help:"Search airports"`
	Airlines  AirlinesCmd  `cmd:"" help:"Look up airlines"`
	Schema    SchemaCmd    `cmd:"" help:"Show CLI schema for agent introspection"`
	ExitCodes ExitCodesCmd `cmd:"" name:"exit-codes" help:"Show exit code reference"`
}

type Globals struct {
	Ctx     context.Context
	Config  *config.Config
	OutOpts output.Options
	CLI     *CLI
	Version string
}

func NewGlobals(cli *CLI) (*Globals, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, errfmt.Config("failed to load config: " + err.Error())
	}

	mode := output.ModeHuman
	if cli.JSON {
		mode = output.ModeJSON
	} else if cli.Plain {
		mode = output.ModePlain
	} else if os.Getenv("AMADEUS_AUTO_JSON") == "1" && !output.IsTerminal() {
		mode = output.ModeJSON
	}

	opts := output.Options{
		Mode:        mode,
		ResultsOnly: cli.ResultsOnly,
		Select:      cli.Select,
		Pretty:      true,
	}

	ctx := output.WithOptions(context.Background(), opts)

	return &Globals{
		Ctx:    ctx,
		Config: cfg,
		OutOpts: opts,
		CLI:    cli,
	}, nil
}

// NewAPIClient creates an authenticated API client, loading and refreshing the token as needed.
func (g *Globals) NewAPIClient() (*api.Client, error) {
	env := g.Config.ResolveEnvironment(g.CLI.Test)
	apiKey := g.Config.ResolveAPIKey()
	apiSecret := g.Config.ResolveAPISecret()

	tok, err := auth.LoadToken(env)
	if err != nil || tok.IsExpired() {
		if apiKey == "" || apiSecret == "" {
			return nil, errfmt.Auth("not authenticated — run 'amadeus auth login' or set AMADEUS_API_KEY and AMADEUS_API_SECRET")
		}
		tok, err = auth.FetchToken(apiKey, apiSecret, env)
		if err != nil {
			return nil, errfmt.Wrap(errfmt.ExitAuth, "authentication failed", err)
		}
		if err := auth.StoreToken(env, tok); err != nil {
			output.Warn("failed to cache token: %s", err)
		}
	}

	return api.NewClient(tok, env), nil
}
```

**Step 2: Update main.go with Kong**

```go
// cmd/amadeus/main.go
package main

import (
	"errors"
	"os"

	"github.com/alecthomas/kong"
	"github.com/voska/amadeus-cli/internal/cmd"
	"github.com/voska/amadeus-cli/internal/errfmt"
	"github.com/voska/amadeus-cli/internal/output"
)

var version = "dev"

func main() {
	var cli cmd.CLI
	ctx := kong.Parse(&cli,
		kong.Name("amadeus"),
		kong.Description("CLI for Amadeus travel APIs"),
		kong.UsageOnError(),
		kong.Vars{"version": version},
	)

	globals, err := cmd.NewGlobals(&cli)
	if err != nil {
		handleError(err)
	}
	globals.Version = version

	if err := ctx.Run(globals); err != nil {
		handleError(err)
	}
}

func handleError(err error) {
	var e *errfmt.Error
	if errors.As(err, &e) {
		output.ErrorMsg("%s", e.Error())
		os.Exit(e.Code)
	}
	output.ErrorMsg("%s", err.Error())
	os.Exit(errfmt.ExitError)
}
```

**Step 3: Add Kong dependency**

Run: `go get github.com/alecthomas/kong@v1.14.0`

**Step 4: Create stub commands** (just enough for compilation)

Create stub files for `AuthCmd`, `FlightsCmd`, `HotelsCmd`, `AirportsCmd`, `AirlinesCmd`, `SchemaCmd` — each as an empty struct with a `Run(g *Globals) error` method that returns nil.

**Step 5: Verify build and help text**

Run: `make build && ./bin/amadeus --help`
Expected: Shows all commands and global flags

**Step 6: Commit**

```bash
git add cmd/ internal/cmd/root.go go.mod go.sum && git commit -m "feat: Kong CLI wiring with global flags and command stubs"
```

---

### Task 9: Auth commands

**Files:**
- Modify: `internal/cmd/auth.go`

**Step 1: Implement auth login, status, logout**

```go
// internal/cmd/auth.go
package cmd

import (
	"github.com/voska/amadeus-cli/internal/auth"
	"github.com/voska/amadeus-cli/internal/errfmt"
	"github.com/voska/amadeus-cli/internal/output"
)

type AuthCmd struct {
	Login  AuthLoginCmd  `cmd:"" help:"Authenticate with Amadeus API"`
	Status AuthStatusCmd `cmd:"" help:"Show authentication status"`
	Logout AuthLogoutCmd `cmd:"" help:"Remove stored credentials"`
}

type AuthLoginCmd struct{}

func (c *AuthLoginCmd) Run(g *Globals) error {
	env := g.Config.ResolveEnvironment(g.CLI.Test)
	apiKey := g.Config.ResolveAPIKey()
	apiSecret := g.Config.ResolveAPISecret()

	if apiKey == "" || apiSecret == "" {
		return errfmt.Config("AMADEUS_API_KEY and AMADEUS_API_SECRET must be set (env vars or ~/.config/amadeus/config.json)")
	}

	output.Hint("authenticating with %s environment...", env)
	tok, err := auth.FetchToken(apiKey, apiSecret, env)
	if err != nil {
		return errfmt.Wrap(errfmt.ExitAuth, "login failed", err)
	}

	if err := auth.StoreToken(env, tok); err != nil {
		return errfmt.Wrap(errfmt.ExitError, "failed to store token", err)
	}

	output.Success("authenticated (%s)", env)
	return nil
}

type AuthStatusCmd struct{}

func (c *AuthStatusCmd) Run(g *Globals) error {
	env := g.Config.ResolveEnvironment(g.CLI.Test)
	tok, err := auth.LoadToken(env)
	if err != nil {
		return errfmt.Auth("not authenticated — run 'amadeus auth login'")
	}

	status := map[string]any{
		"authenticated": true,
		"environment":   env,
		"expired":       tok.IsExpired(),
		"expires_at":    tok.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return output.Write(g.Ctx, status)
}

type AuthLogoutCmd struct{}

func (c *AuthLogoutCmd) Run(g *Globals) error {
	env := g.Config.ResolveEnvironment(g.CLI.Test)
	if err := auth.DeleteToken(env); err != nil {
		output.Warn("no stored token for %s", env)
		return nil
	}
	output.Success("logged out (%s)", env)
	return nil
}
```

**Step 2: Test manually**

Run: `make build && ./bin/amadeus auth status`
Expected: Error message about not being authenticated (exit code 4)

**Step 3: Commit**

```bash
git add internal/cmd/auth.go && git commit -m "feat: auth login, status, logout commands"
```

---

### Task 10: Flight search command

**Files:**
- Modify: `internal/cmd/flights.go`

**Step 1: Implement flights search**

```go
// internal/cmd/flights.go
package cmd

import (
	"net/url"
	"strconv"

	"github.com/voska/amadeus-cli/internal/errfmt"
	"github.com/voska/amadeus-cli/internal/output"
)

type FlightsCmd struct {
	Search  FlightsSearchCmd  `cmd:"" help:"Search for flight offers"`
	Price   FlightsPriceCmd   `cmd:"" help:"Confirm pricing for a flight offer"`
	Seatmap FlightsSeatmapCmd `cmd:"" help:"View seat map for a flight offer"`
}

type FlightsSearchCmd struct {
	From        string `required:"" name:"from" help:"Origin IATA code (e.g., JFK)"`
	To          string `required:"" name:"to" help:"Destination IATA code (e.g., CDG)"`
	Date        string `required:"" name:"date" help:"Departure date (YYYY-MM-DD)"`
	Return      string `name:"return" help:"Return date for round-trip (YYYY-MM-DD)"`
	Adults      int    `name:"adults" default:"1" help:"Number of adults (1-9)"`
	Children    int    `name:"children" default:"0" help:"Number of children (0-9)"`
	TravelClass string `name:"class" help:"Travel class (ECONOMY, PREMIUM_ECONOMY, BUSINESS, FIRST)"`
	NonStop     bool   `name:"nonstop" help:"Direct flights only"`
	Currency    string `name:"currency" help:"Currency code (e.g., USD, EUR)"`
	MaxPrice    int    `name:"max-price" help:"Maximum price (no decimals)"`
	Max         int    `name:"max" default:"10" help:"Maximum number of results (1-250)"`
}

func (c *FlightsSearchCmd) Run(g *Globals) error {
	client, err := g.NewAPIClient()
	if err != nil {
		return err
	}

	params := url.Values{
		"originLocationCode":      {c.From},
		"destinationLocationCode": {c.To},
		"departureDate":           {c.Date},
		"adults":                  {strconv.Itoa(c.Adults)},
		"max":                     {strconv.Itoa(c.Max)},
	}
	if c.Return != "" {
		params.Set("returnDate", c.Return)
	}
	if c.Children > 0 {
		params.Set("children", strconv.Itoa(c.Children))
	}
	if c.TravelClass != "" {
		params.Set("travelClass", c.TravelClass)
	}
	if c.NonStop {
		params.Set("nonStop", "true")
	}
	if c.Currency != "" {
		params.Set("currencyCode", c.Currency)
	}
	if c.MaxPrice > 0 {
		params.Set("maxPrice", strconv.Itoa(c.MaxPrice))
	}

	if g.CLI.DryRun {
		output.Hint("[dry-run] GET /v2/shopping/flight-offers?%s", params.Encode())
		return nil
	}

	result, err := client.Get("/v2/shopping/flight-offers", params)
	if err != nil {
		return err
	}

	data, ok := result["data"].([]any)
	if !ok || len(data) == 0 {
		return errfmt.Empty()
	}

	return output.Write(g.Ctx, result)
}
```

Add stub `FlightsPriceCmd` and `FlightsSeatmapCmd` that return "not yet implemented" for now.

**Step 2: Verify build**

Run: `make build && ./bin/amadeus flights search --help`
Expected: Shows all flight search flags

**Step 3: Commit**

```bash
git add internal/cmd/flights.go && git commit -m "feat: flights search command with all Amadeus query parameters"
```

---

### Task 11: Hotel commands

**Files:**
- Modify: `internal/cmd/hotels.go`

**Step 1: Implement hotels search, offers, book**

```go
// internal/cmd/hotels.go
package cmd

type HotelsCmd struct {
	Search HotelsSearchCmd `cmd:"" help:"Search hotels by city or location"`
	Offers HotelsOffersCmd `cmd:"" help:"Get offers for a specific hotel"`
	Book   HotelsBookCmd   `cmd:"" help:"Book a hotel offer"`
}

type HotelsSearchCmd struct {
	City    string  `name:"city" help:"City IATA code (e.g., PAR)" xor:"location"`
	Lat     float64 `name:"lat" help:"Latitude" xor:"location"`
	Lng     float64 `name:"lng" help:"Longitude" xor:"location"`
	Radius  int     `name:"radius" default:"5" help:"Search radius in km"`
	Ratings []int   `name:"ratings" help:"Hotel ratings to filter (1-5)"`
}

// Run: GET /v1/reference-data/locations/hotels/by-city or by-geocode

type HotelsOffersCmd struct {
	HotelID  string `required:"" name:"hotel-id" help:"Amadeus hotel ID"`
	CheckIn  string `required:"" name:"checkin" help:"Check-in date (YYYY-MM-DD)"`
	CheckOut string `required:"" name:"checkout" help:"Check-out date (YYYY-MM-DD)"`
	Adults   int    `name:"adults" default:"1" help:"Number of adults"`
	Rooms    int    `name:"rooms" default:"1" help:"Number of rooms"`
	Currency string `name:"currency" help:"Currency code"`
}

// Run: GET /v3/shopping/hotel-offers?hotelIds=XXX

type HotelsBookCmd struct {
	OfferID    string `required:"" name:"offer-id" help:"Hotel offer ID from search results"`
	GuestName  string `required:"" name:"guest-name" help:"Guest full name"`
	GuestEmail string `required:"" name:"guest-email" help:"Guest email address"`
}

// Run: POST /v2/booking/hotel-orders
```

Implement `Run(g *Globals) error` for each, following the same pattern as flights: build params, call client.Get/Post, handle empty results, write output.

**Step 2: Verify build**

Run: `make build && ./bin/amadeus hotels search --help`

**Step 3: Commit**

```bash
git add internal/cmd/hotels.go && git commit -m "feat: hotel search, offers, and booking commands"
```

---

### Task 12: Airport and airline lookup commands

**Files:**
- Modify: `internal/cmd/airports.go`
- Modify: `internal/cmd/airlines.go`

**Step 1: Implement airport search**

```go
type AirportsCmd struct {
	Search AirportsSearchCmd `cmd:"" help:"Search airports by keyword (autocomplete)"`
}

type AirportsSearchCmd struct {
	Keyword string `required:"" arg:"" help:"Search keyword (e.g., 'london', 'JFK')"`
}

// Run: GET /v1/reference-data/locations?subType=AIRPORT&keyword=XXX
```

**Step 2: Implement airline lookup**

```go
type AirlinesCmd struct {
	Lookup AirlinesLookupCmd `cmd:"" name:"lookup" help:"Look up airline by IATA code"`
}

type AirlinesLookupCmd struct {
	Code string `required:"" arg:"" help:"IATA airline code (e.g., BA, AA)"`
}

// Run: GET /v1/reference-data/airlines?airlineCodes=XX
```

**Step 3: Verify build**

Run: `make build && ./bin/amadeus airports search --help`

**Step 4: Commit**

```bash
git add internal/cmd/airports.go internal/cmd/airlines.go && git commit -m "feat: airport search and airline lookup commands"
```

---

### Task 13: Schema introspection command

**Files:**
- Modify: `internal/cmd/schema.go`

**Step 1: Implement schema command**

Reference: `/Users/matt/Developer/qbo-cli/internal/cmd/schema.go`

Build a `fullSchema(version string) map[string]any` that returns the complete CLI tree: all commands, subcommands, flags, and their types. This is the primary agent introspection endpoint.

Also wire `ExitCodesCmd.Run` to output the exit code table.

**Step 2: Verify output**

Run: `make build && ./bin/amadeus schema --json | jq .`
Expected: Full CLI tree as JSON

Run: `./bin/amadeus exit-codes --json | jq .`
Expected: Exit code table as JSON

**Step 3: Commit**

```bash
git add internal/cmd/schema.go && git commit -m "feat: schema introspection and exit-codes commands for agent use"
```

---

### Task 14: GoReleaser and distribution

**Files:**
- Create: `.goreleaser.yaml`

**Step 1: Create GoReleaser config**

Reference: `/Users/matt/Developer/qbo-cli/.goreleaser.yaml`

```yaml
version: 2
builds:
  - main: ./cmd/amadeus
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w -X main.version={{.Version}}

archives:
  - format: tar.gz
    name_template: "amadeus_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: "checksums.txt"

brews:
  - repository:
      owner: voska
      name: homebrew-tap
    name: amadeus
    homepage: "https://github.com/voska/amadeus-cli"
    description: "CLI for Amadeus travel APIs"
    install: |
      bin.install "amadeus"
```

**Step 2: Verify GoReleaser config**

Run: `goreleaser check` (if installed)

**Step 3: Commit**

```bash
git add .goreleaser.yaml && git commit -m "feat: GoReleaser config for cross-platform distribution"
```

---

### Task 15: Claude Code skill

**Files:**
- Create: `skills/amadeus/SKILL.md`
- Create: `skills/amadeus/references/COMMANDS.md`

**Step 1: Create agent skill file**

Write a SKILL.md with: overview, auth setup, command examples, common workflows (search flights, compare prices, book hotel), and tips for agents.

**Step 2: Create command reference**

Auto-generate from `amadeus schema --json` or write manually with all commands, flags, and examples.

**Step 3: Commit**

```bash
git add skills/ && git commit -m "feat: Claude Code skill for agent integration"
```

---

### Task 16: End-to-end testing

**Step 1: Get Amadeus test API credentials**

Sign up at https://developers.amadeus.com/ and get test API key + secret.

**Step 2: Test auth flow**

```bash
export AMADEUS_API_KEY=your_key
export AMADEUS_API_SECRET=your_secret
./bin/amadeus auth login --test
./bin/amadeus auth status --test --json
```

**Step 3: Test flight search**

```bash
./bin/amadeus flights search --from JFK --to CDG --date 2026-04-15 --test --json | jq '.data[0].price'
./bin/amadeus flights search --from JFK --to CDG --date 2026-04-15 --test --results-only --select id,price.total,price.currency
```

**Step 4: Test hotel search**

```bash
./bin/amadeus hotels search --city PAR --test --json
./bin/amadeus hotels offers --hotel-id XXXX --checkin 2026-04-15 --checkout 2026-04-18 --test --json
```

**Step 5: Test agent features**

```bash
./bin/amadeus schema --json | jq '.commands[].name'
./bin/amadeus exit-codes --json
echo "pipe test:" && AMADEUS_AUTO_JSON=1 ./bin/amadeus flights search --from JFK --to CDG --date 2026-04-15 --test | jq .
```

**Step 6: Fix any issues found during testing**

**Step 7: Commit any fixes**

```bash
git add -A && git commit -m "fix: adjustments from end-to-end testing"
```
