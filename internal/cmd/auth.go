package cmd

type AuthCmd struct {
	Login  AuthLoginCmd  `cmd:"" help:"Authenticate with Amadeus API"`
	Status AuthStatusCmd `cmd:"" help:"Show authentication status"`
	Logout AuthLogoutCmd `cmd:"" help:"Remove stored credentials"`
}

type AuthLoginCmd struct{}

func (c *AuthLoginCmd) Run(g *Globals) error { return nil }

type AuthStatusCmd struct{}

func (c *AuthStatusCmd) Run(g *Globals) error { return nil }

type AuthLogoutCmd struct{}

func (c *AuthLogoutCmd) Run(g *Globals) error { return nil }
