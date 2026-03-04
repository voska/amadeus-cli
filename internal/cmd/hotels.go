package cmd

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/voska/amadeus-cli/internal/errfmt"
	"github.com/voska/amadeus-cli/internal/output"
)

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

func (c *HotelsSearchCmd) Run(g *Globals) error {
	var path string
	params := url.Values{}

	if c.City != "" {
		path = "/v1/reference-data/locations/hotels/by-city"
		params.Set("cityCode", c.City)
	} else if c.Lat != 0 || c.Lng != 0 {
		path = "/v1/reference-data/locations/hotels/by-geocode"
		params.Set("latitude", fmt.Sprintf("%.4f", c.Lat))
		params.Set("longitude", fmt.Sprintf("%.4f", c.Lng))
	} else {
		return errfmt.Usage("provide --city or --lat/--lng")
	}

	if c.Radius > 0 {
		params.Set("radius", strconv.Itoa(c.Radius))
		params.Set("radiusUnit", "KM")
	}
	if len(c.Ratings) > 0 {
		ratings := make([]string, len(c.Ratings))
		for i, r := range c.Ratings {
			ratings[i] = strconv.Itoa(r)
		}
		params.Set("ratings", strings.Join(ratings, ","))
	}

	if g.CLI.DryRun {
		output.Hint("[dry-run] GET %s?%s", path, params.Encode())
		return nil
	}

	client, err := g.NewAPIClient()
	if err != nil {
		return err
	}

	result, err := client.Get(path, params)
	if err != nil {
		return err
	}

	data, ok := result["data"].([]any)
	if !ok || len(data) == 0 {
		return errfmt.Empty()
	}

	return output.Write(g.Ctx, result)
}

type HotelsOffersCmd struct {
	HotelID  string `required:"" name:"hotel-id" help:"Amadeus hotel ID"`
	CheckIn  string `required:"" name:"checkin" help:"Check-in date (YYYY-MM-DD)"`
	CheckOut string `required:"" name:"checkout" help:"Check-out date (YYYY-MM-DD)"`
	Adults   int    `name:"adults" default:"1" help:"Number of adults"`
	Rooms    int    `name:"rooms" default:"1" help:"Number of rooms"`
	Currency string `name:"currency" help:"Currency code"`
}

func (c *HotelsOffersCmd) Run(g *Globals) error {
	params := url.Values{
		"hotelIds":     {c.HotelID},
		"checkInDate":  {c.CheckIn},
		"checkOutDate": {c.CheckOut},
		"adults":       {strconv.Itoa(c.Adults)},
		"roomQuantity": {strconv.Itoa(c.Rooms)},
	}
	if c.Currency != "" {
		params.Set("currency", c.Currency)
	}

	if g.CLI.DryRun {
		output.Hint("[dry-run] GET /v3/shopping/hotel-offers?%s", params.Encode())
		return nil
	}

	client, err := g.NewAPIClient()
	if err != nil {
		return err
	}

	result, err := client.Get("/v3/shopping/hotel-offers", params)
	if err != nil {
		return err
	}

	data, ok := result["data"].([]any)
	if !ok || len(data) == 0 {
		return errfmt.Empty()
	}

	return output.Write(g.Ctx, result)
}

type HotelsBookCmd struct {
	OfferID    string `required:"" name:"offer-id" help:"Hotel offer ID from search results"`
	GuestName  string `required:"" name:"guest-name" help:"Guest full name"`
	GuestEmail string `required:"" name:"guest-email" help:"Guest email address"`
}

func (c *HotelsBookCmd) Run(g *Globals) error {
	body := map[string]any{
		"data": map[string]any{
			"type":    "hotel-order",
			"offerId": c.OfferID,
			"guests": []map[string]any{
				{
					"tid": 1,
					"name": map[string]any{
						"firstName": strings.Split(c.GuestName, " ")[0],
						"lastName":  lastNameFromFull(c.GuestName),
					},
					"contact": map[string]any{
						"email": c.GuestEmail,
					},
				},
			},
		},
	}

	if g.CLI.DryRun {
		output.Hint("[dry-run] POST /v2/booking/hotel-orders")
		return nil
	}

	client, err := g.NewAPIClient()
	if err != nil {
		return err
	}

	result, err := client.Post("/v2/booking/hotel-orders", body)
	if err != nil {
		return err
	}

	return output.Write(g.Ctx, result)
}

func lastNameFromFull(full string) string {
	parts := strings.Fields(full)
	if len(parts) > 1 {
		return strings.Join(parts[1:], " ")
	}
	return full
}
