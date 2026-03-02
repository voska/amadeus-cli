package cmd

type AirportsCmd struct {
	Search AirportsSearchCmd `cmd:"" help:"Search airports by keyword (autocomplete)"`
}

type AirportsSearchCmd struct{}

func (c *AirportsSearchCmd) Run(g *Globals) error { return nil }
