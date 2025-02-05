package config_mutable_cli

import (
	"flag"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/bravo/options_tools"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
)

type Config struct {
	BasePath string

	Debug            debug.Options
	Verbose          bool
	Quiet            bool
	Todo             bool
	dryRun           bool
	Complete         bool
	IgnoreHookErrors bool
	Hooks            string

	CheckoutCacheEnabled bool
	PredictableZettelIds bool

	PrintOptions, maskPrintOptions options_print.V0
	ToolOptions                    options_tools.Options

	descriptions.Description
}

// TODO add support for all flags
func (c Config) GetCLIFlags() (flags []string) {
	flags = append(flags, fmt.Sprintf("-print-time=%t", c.PrintOptions.PrintTime))
	flags = append(flags, fmt.Sprintf("-print-colors=%t", c.PrintOptions.PrintColors))

	if c.Verbose {
		flags = append(flags, "-verbose")
	}

	return
}

func (c *Config) SetFlagSet(f *flag.FlagSet) {
	f.StringVar(&c.BasePath, "dir-zit", "", "")

	f.Var(&c.Debug, "debug", "debugging options")
	f.BoolVar(&c.Todo, "todo", false, "")
	f.BoolVar(&c.dryRun, "dry-run", false, "")
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

func Default() (c Config) {
	c.PrintOptions = options_print.Default()

	return
}

func (c *Config) ApplyPrintOptionsConfig(
	po options_print.V0,
) {
	cliSet := c.PrintOptions
	c.PrintOptions = po
	c.PrintOptions.Merge(cliSet, c.maskPrintOptions)
}

func (c Config) UsePredictableZettelIds() bool {
	return c.PredictableZettelIds
}

func (c Config) UsePrintTime() bool {
	return c.PrintOptions.PrintTime
}

func (c Config) UsePrintTags() bool {
	return c.PrintOptions.PrintTagsAlways
}

func (c Config) IsDryRun() bool {
	return c.dryRun || c.Complete
}

func (c *Config) SetDryRun(v bool) {
	c.dryRun = v
}
