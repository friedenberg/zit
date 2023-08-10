package erworben

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/erworben_cli_print_options"
	"github.com/friedenberg/zit/src/charlie/debug"
)

type Cli struct {
	BasePath string

	Debug    debug.Options
	Verbose  bool
	Todo     bool
	DryRun   bool
	Complete bool

	IncludeCwd    bool
	IncludeHidden bool

	AllowMissingHinweis              bool
	CheckoutCacheEnabled             bool
	PredictableHinweisen             bool
	UseRightAlignedIndentsInOrganize bool

	erworben_cli_print_options.Options
}

func (c *Cli) AddToFlags(f *flag.FlagSet) {
	f.StringVar(&c.BasePath, "dir-zit", "", "")

	f.Var(&c.Debug, "debug", "debugging options")
	f.BoolVar(&c.Todo, "todo", false, "")
	f.BoolVar(&c.DryRun, "dry-run", false, "")
	f.BoolVar(&c.Verbose, "verbose", false, "")
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

	c.Options.AddToFlags(f)
}

func DefaultCli() (c Cli) {
	return
}

func (c Cli) UsePredictableHinweisen() bool {
	return c.PredictableHinweisen
}

func (c Cli) UsePrintTime() bool {
	return c.PrintTime
}

func (c Cli) UsePrintEtiketten() bool {
	return c.PrintEtikettenAlways
}
