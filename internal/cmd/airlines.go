package cmd

import (
	"net/url"

	"github.com/voska/amadeus-cli/internal/errfmt"
	"github.com/voska/amadeus-cli/internal/output"
)

type AirlinesCmd struct {
	Lookup AirlinesLookupCmd `cmd:"" name:"lookup" help:"Look up airline by IATA code"`
}

type AirlinesLookupCmd struct {
	Code string `required:"" arg:"" help:"IATA airline code (e.g., BA, AA)"`
}

func (c *AirlinesLookupCmd) Run(g *Globals) error {
	params := url.Values{
		"airlineCodes": {c.Code},
	}

	if g.CLI.DryRun {
		output.Hint("[dry-run] GET /v1/reference-data/airlines?%s", params.Encode())
		return nil
	}

	client, err := g.NewAPIClient()
	if err != nil {
		return err
	}

	result, err := client.Get("/v1/reference-data/airlines", params)
	if err != nil {
		return err
	}

	data, ok := result["data"].([]any)
	if !ok || len(data) == 0 {
		return errfmt.Empty()
	}

	return output.Write(g.Ctx, result)
}
