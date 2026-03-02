package cmd

type AirlinesCmd struct {
	Lookup AirlinesLookupCmd `cmd:"" name:"lookup" help:"Look up airline by IATA code"`
}

type AirlinesLookupCmd struct{}

func (c *AirlinesLookupCmd) Run(g *Globals) error { return nil }
