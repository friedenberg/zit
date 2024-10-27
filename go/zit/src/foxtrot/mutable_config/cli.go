package mutable_config

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/bravo/erworben_tools"
	"code.linenisgreat.com/zit/go/zit/src/charlie/print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
)

type Cli struct {
	BasePath string

	Debug            debug.Options
	Verbose          bool
	Quiet            bool
	Todo             bool
	DryRun           bool
	Complete         bool
	IgnoreHookErrors bool
	Hooks            string

	CheckoutCacheEnabled bool
	PredictableZettelIds bool

	PrintOptions, maskPrintOptions print_options.General
	ToolOptions                    erworben_tools.Tools

	descriptions.Description
}

func (c *Cli) AddToFlags(f *flag.FlagSet) {
	f.StringVar(&c.BasePath, "dir-zit", "", "")

	f.Var(&c.Debug, "debug", "debugging options")
	f.BoolVar(&c.Todo, "todo", false, "")
	f.BoolVar(&c.DryRun, "dry-run", false, "")
	f.BoolVar(&c.Verbose, "verbose", false, "")
	f.BoolVar(&c.Quiet, "quiet", false, "")
	f.BoolVar(&c.Complete, "complete", false, "")

	f.BoolVar(&c.CheckoutCacheEnabled, "checkout-cache-enabled", false, "")

	f.BoolVar(
		&c.PredictableZettelIds,
		"predictable-zettel-ids",
		false,
		"generate new zettel ids in order",
	)

	c.PrintOptions.AddToFlags(f, &c.maskPrintOptions)
	c.ToolOptions.AddToFlags(f)

	f.BoolVar(
		&c.PrintOptions.ZittishNewlines,
		"zittish-newlines",
		false,
		"add extra newlines to zittish to improve readability",
	)

	f.BoolVar(
		&c.IgnoreHookErrors,
		"ignore-hook-errors",
		false,
		"ignores errors coming out of hooks",
	)

	f.StringVar(&c.Hooks, "hooks", "", "")

	f.Var(&c.Description, "komment", "Comment for Bestandsaufnahme")
}

func DefaultCli() (c Cli) {
	c.PrintOptions = print_options.Default()

	return
}

func (c *Cli) ApplyPrintOptionsConfig(
	po print_options.General,
) {
	cliSet := c.PrintOptions
	c.PrintOptions = po
	c.PrintOptions.Merge(cliSet, c.maskPrintOptions)
}

func (c Cli) UsePredictableHinweisen() bool {
	return c.PredictableZettelIds
}

func (c Cli) UsePrintTime() bool {
	return c.PrintOptions.PrintTime
}

func (c Cli) UsePrintTags() bool {
	return c.PrintOptions.PrintTagsAlways
}
