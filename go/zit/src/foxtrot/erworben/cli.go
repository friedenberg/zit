package erworben

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/bravo/erworben_tools"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/echo/bezeichnung"
)

type Cli struct {
	BasePath string

	Debug         debug.Options
	Verbose       bool
	Quiet         bool
	Todo          bool
	DryRun        bool
	Complete      bool
	ChrestEnabled bool
	Hooks         string

	IncludeCwd    bool
	IncludeHidden bool

	AllowMissingHinweis  bool
	CheckoutCacheEnabled bool
	PredictableHinweisen bool
	EtikettenPaths       bool

	PrintOptions, maskPrintOptions erworben_cli_print_options.PrintOptions
	ToolOptions                    erworben_tools.Tools

	bezeichnung.Bezeichnung
}

func (c *Cli) AddToFlags(f *flag.FlagSet) {
	f.StringVar(&c.BasePath, "dir-zit", "", "")

	f.Var(&c.Debug, "debug", "debugging options")
	f.BoolVar(&c.Todo, "todo", false, "")
	f.BoolVar(&c.DryRun, "dry-run", false, "")
	f.BoolVar(&c.Verbose, "verbose", false, "")
	f.BoolVar(&c.Quiet, "quiet", false, "")
	f.BoolVar(&c.Complete, "complete", false, "")

	f.BoolVar(
		&c.IncludeCwd,
		"include-cwd",
		true,
		"include checked-out Objekten in the working directory",
	)

	f.BoolVar(
		&c.IncludeHidden,
		"include-hidden",
		false,
		"include zettels that have hidden etiketten",
	)

	f.BoolVar(&c.AllowMissingHinweis, "allow-missing-hinweis", false, "")

	f.BoolVar(&c.CheckoutCacheEnabled, "checkout-cache-enabled", false, "")

	f.BoolVar(
		&c.PredictableHinweisen,
		"predictable-hinweisen",
		false,
		"don't randomly select new hinweisen",
	)

	f.BoolVar(
		&c.EtikettenPaths,
		"etiketten-paths",
		true,
		"use new etiketten paths for queries",
	)

	f.BoolVar(
		&c.ChrestEnabled,
		"chrest-enabled",
		true,
		"enable Chrest functionality",
	)

	c.PrintOptions.AddToFlags(f, &c.maskPrintOptions)
	c.ToolOptions.AddToFlags(f)

	f.BoolVar(
		&c.PrintOptions.ZittishNewlines,
		"zittish-newlines",
		false,
		"add extra newlines to zittish to improve readability",
	)

	f.StringVar(&c.Hooks, "hooks", "", "")

	f.Var(&c.Bezeichnung, "komment", "Comment for Bestandsaufnahme")
}

func DefaultCli() (c Cli) {
	c.PrintOptions = erworben_cli_print_options.Default()

	return
}

func (c *Cli) ApplyPrintOptionsKonfig(
	po erworben_cli_print_options.PrintOptions,
) {
	cliSet := c.PrintOptions
	c.PrintOptions = po
	c.PrintOptions.Merge(cliSet, c.maskPrintOptions)
}

func (c Cli) UsePredictableHinweisen() bool {
	return c.PredictableHinweisen
}

func (c Cli) UsePrintTime() bool {
	return c.PrintOptions.PrintTime
}

func (c Cli) UsePrintEtiketten() bool {
	return c.PrintOptions.PrintEtikettenAlways
}
