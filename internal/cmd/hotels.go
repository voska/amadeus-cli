package cmd

type HotelsCmd struct {
	Search HotelsSearchCmd `cmd:"" help:"Search hotels by city or location"`
	Offers HotelsOffersCmd `cmd:"" help:"Get offers for a specific hotel"`
	Book   HotelsBookCmd   `cmd:"" help:"Book a hotel offer"`
}

type HotelsSearchCmd struct{}

func (c *HotelsSearchCmd) Run(g *Globals) error { return nil }

type HotelsOffersCmd struct{}

func (c *HotelsOffersCmd) Run(g *Globals) error { return nil }

type HotelsBookCmd struct{}

func (c *HotelsBookCmd) Run(g *Globals) error { return nil }
