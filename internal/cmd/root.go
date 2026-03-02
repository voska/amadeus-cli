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
	JSON        bool     `short:"j" name:"json" help:"Output as JSON" aliases:"machine"`
	Plain       bool     `short:"p" name:"plain" help:"Output as tab-separated values"`
	ResultsOnly bool     `name:"results-only" help:"Strip response metadata, return data array only"`
	Select      []string `name:"select" short:"s" help:"Select specific fields (dot-path)" aliases:"fields"`
	Test        bool     `name:"test" help:"Use test environment (test.api.amadeus.com)"`
	DryRun      bool     `name:"dry-run" help:"Show what would be done without making API calls"`
	NoInput     bool     `name:"no-input" help:"Never prompt for input"`
	Verbose     bool     `short:"v" name:"verbose" help:"Verbose output"`

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
		Ctx:     ctx,
		Config:  cfg,
		OutOpts: opts,
		CLI:     cli,
	}, nil
}

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
