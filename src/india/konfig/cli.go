package konfig

import (
	"flag"

	"github.com/friedenberg/zit/src/charlie/debug"
)

type Cli struct {
	BasePath string

	Debug    debug.Options
	Verbose  bool
	DryRun   bool
	Complete bool

	IncludeHidden  bool
	IncludeHistory bool

	AllowMissingHinweis              bool
	CheckoutCacheEnabled             bool
	PredictableHinweisen             bool
	UseRightAlignedIndentsInOrganize bool
	PrintAbbreviatedHinweisen        bool
	PrintAbbreviatedShas             bool
	PrintNewShaSyntax                bool
	PrintIncludeTypen                bool
	PrintIncludeBezeichnungen        bool
}

func (c *Cli) AddToFlags(f *flag.FlagSet) {
	f.StringVar(&c.BasePath, "dir-zit", "", "")

	f.Var(&c.Debug, "debug", "debugging options")
	f.BoolVar(&c.Verbose, "verbose", false, "")
	f.BoolVar(&c.DryRun, "dry-run", false, "")
	f.BoolVar(&c.Complete, "complete", false, "")

	f.BoolVar(&c.IncludeHidden, "include-hidden", false, "include zettels that have hidden etiketten")
	f.BoolVar(&c.IncludeHistory, "include-history", false, "include past zettels")

	f.BoolVar(&c.CheckoutCacheEnabled, "checkout-cache-enabled", false, "")
	f.BoolVar(&c.AllowMissingHinweis, "allow-missing-hinweis", false, "")
	f.BoolVar(&c.PredictableHinweisen, "predictable-hinweisen", false, "don't randomly select new hinweisen")

	f.BoolVar(&c.PrintAbbreviatedShas, "abbreviate-shas", true, "")
	f.BoolVar(&c.PrintAbbreviatedHinweisen, "abbreviate-hinweisen", true, "")
	f.BoolVar(&c.PrintNewShaSyntax, "new-zettel-sha-syntax", true, "")
	f.BoolVar(&c.PrintIncludeTypen, "print-typen", true, "")
	f.BoolVar(&c.PrintIncludeBezeichnungen, "print-bezeichnungen", true, "")
}

func DefaultCli() (c Cli) {
	return
}
