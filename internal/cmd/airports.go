package cmd

import (
	"net/url"

	"github.com/voska/amadeus-cli/internal/errfmt"
	"github.com/voska/amadeus-cli/internal/output"
)

type AirportsCmd struct {
	Search AirportsSearchCmd `cmd:"" help:"Search airports by keyword (autocomplete)"`
}

type AirportsSearchCmd struct {
	Keyword string `required:"" arg:"" help:"Search keyword (e.g., 'london', 'JFK')"`
}

func (c *AirportsSearchCmd) Run(g *Globals) error {
	params := url.Values{
		"subType": {"AIRPORT"},
		"keyword": {c.Keyword},
	}

	if g.CLI.DryRun {
		output.Hint("[dry-run] GET /v1/reference-data/locations?%s", params.Encode())
		return nil
	}

	client, err := g.NewAPIClient()
	if err != nil {
		return err
	}

	result, err := client.Get("/v1/reference-data/locations", params)
	if err != nil {
		return err
	}

	data, ok := result["data"].([]any)
	if !ok || len(data) == 0 {
		return errfmt.Empty()
	}

	return output.Write(g.Ctx, result)
}
