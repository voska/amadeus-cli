package cmd

import (
	"github.com/voska/amadeus-cli/internal/errfmt"
	"github.com/voska/amadeus-cli/internal/output"
)

type ExitCodeEntry struct {
	Code        int    `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ExitCodesCmd struct{}

func (c *ExitCodesCmd) Run(g *Globals) error {
	table := ExitCodeTable()
	entries := make([]any, len(table))
	for i, e := range table {
		entries[i] = map[string]any{
			"code":        e.Code,
			"name":        e.Name,
			"description": e.Description,
		}
	}
	return output.Write(g.Ctx, entries)
}

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
