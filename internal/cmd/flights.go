package cmd

type FlightsCmd struct {
	Search  FlightsSearchCmd  `cmd:"" help:"Search for flight offers"`
	Price   FlightsPriceCmd   `cmd:"" help:"Confirm pricing for a flight offer"`
	Seatmap FlightsSeatmapCmd `cmd:"" help:"View seat map for a flight offer"`
}

type FlightsSearchCmd struct{}

func (c *FlightsSearchCmd) Run(g *Globals) error { return nil }

type FlightsPriceCmd struct{}

func (c *FlightsPriceCmd) Run(g *Globals) error { return nil }

type FlightsSeatmapCmd struct{}

func (c *FlightsSeatmapCmd) Run(g *Globals) error { return nil }
