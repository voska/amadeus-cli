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
