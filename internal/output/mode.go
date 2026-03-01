package output

import (
	"context"
	"os"
)

type Mode int

const (
	ModeHuman Mode = iota
	ModeJSON
	ModePlain
)

type Options struct {
	Mode        Mode
	ResultsOnly bool
	Select      []string
	Pretty      bool
}

type ctxKey struct{}

func WithOptions(ctx context.Context, opts Options) context.Context {
	return context.WithValue(ctx, ctxKey{}, opts)
}

func GetOptions(ctx context.Context) Options {
	if opts, ok := ctx.Value(ctxKey{}).(Options); ok {
		return opts
	}
	return Options{Mode: ModeHuman, Pretty: true}
}

func IsTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
