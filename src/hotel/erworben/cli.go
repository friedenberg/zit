package erworben

import (
	"flag"

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
	PrintAbbreviatedHinweisen        bool
	PrintAbbreviatedKennungen        bool
	PrintAbbreviatedShas             bool
	PrintNewShaSyntax                bool
	PrintIncludeTypen                bool
	PrintIncludeBezeichnungen        bool
	PrintTime                        bool
}

func (c *Cli) AddToFlags(f *flag.FlagSet) {
	f.StringVar(&c.BasePath, "dir-zit", "", "")

	f.Var(&c.Debug, "debug", "debugging options")
	f.BoolVar(&c.Todo, "todo", false, "")
	f.BoolVar(&c.DryRun, "dry-run", false, "")
	f.BoolVar(&c.Verbose, "verbose", false, "")
	f.BoolVar(&c.Complete, "complete", false, "")

	f.BoolVar(&c.IncludeCwd, "include-cwd", true, "include checked-out Objekten in the working directory")
	f.BoolVar(&c.IncludeHidden, "include-hidden", false, "include zettels that have hidden etiketten")

	f.BoolVar(&c.AllowMissingHinweis, "allow-missing-hinweis", false, "")
	f.BoolVar(&c.CheckoutCacheEnabled, "checkout-cache-enabled", false, "")
	f.BoolVar(&c.PredictableHinweisen, "predictable-hinweisen", false, "don't randomly select new hinweisen")

	f.BoolVar(&c.PrintNewShaSyntax, "new-zettel-sha-syntax", true, "")
	f.BoolVar(&c.PrintIncludeTypen, "print-typen", true, "")
	f.BoolVar(&c.PrintAbbreviatedShas, "abbreviate-shas", true, "")
	f.BoolVar(&c.PrintAbbreviatedHinweisen, "abbreviate-hinweisen", true, "")
	f.BoolVar(&c.PrintAbbreviatedKennungen, "abbreviate-kennungen", true, "")
	f.BoolVar(&c.PrintIncludeBezeichnungen, "print-bezeichnungen", true, "")
	f.BoolVar(&c.PrintTime, "print-time", true, "")
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
