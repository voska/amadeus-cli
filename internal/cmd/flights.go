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

	client, err := g.NewAPIClient()
	if err != nil {
		return err
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

type FlightsPriceCmd struct {
	OfferID string `required:"" name:"offer-id" help:"Flight offer ID from search results"`
}

func (c *FlightsPriceCmd) Run(g *Globals) error {
	return errfmt.New(errfmt.ExitError, "not yet implemented — coming soon")
}

type FlightsSeatmapCmd struct {
	OfferID string `required:"" name:"offer-id" help:"Flight offer ID from search results"`
}

func (c *FlightsSeatmapCmd) Run(g *Globals) error {
	return errfmt.New(errfmt.ExitError, "not yet implemented — coming soon")
}
