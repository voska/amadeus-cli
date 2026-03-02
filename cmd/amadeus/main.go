package main

import (
	"errors"
	"os"

	"github.com/alecthomas/kong"
	"github.com/voska/amadeus-cli/internal/cmd"
	"github.com/voska/amadeus-cli/internal/errfmt"
	"github.com/voska/amadeus-cli/internal/output"
)

var version = "dev"

func main() {
	var cli cmd.CLI
	ctx := kong.Parse(&cli,
		kong.Name("amadeus"),
		kong.Description("CLI for Amadeus travel APIs"),
		kong.UsageOnError(),
		kong.Vars{"version": version},
	)

	globals, err := cmd.NewGlobals(&cli)
	if err != nil {
		handleError(err)
	}
	globals.Version = version

	if err := ctx.Run(globals); err != nil {
		handleError(err)
	}
}

func handleError(err error) {
	var e *errfmt.Error
	if errors.As(err, &e) {
		output.ErrorMsg("%s", e.Error())
		os.Exit(e.Code)
	}
	output.ErrorMsg("%s", err.Error())
	os.Exit(errfmt.ExitError)
}
