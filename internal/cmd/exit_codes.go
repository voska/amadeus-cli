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
