package store_fs

import (
	"flag"

	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/cwd"
)

type OptionsReadExternal struct {
	Zettelen map[kennung.Hinweis]zettel.Transacted
}

type CheckoutOptions struct {
	Cwd          cwd.CwdFiles
	Force        bool
	CheckoutMode objekte.CheckoutMode
	Zettelen     map[kennung.Hinweis]zettel.Transacted
}

func (c *CheckoutOptions) AddToFlagSet(f *flag.FlagSet) {
	f.Var(&c.CheckoutMode, "mode", "mode for checking out the zettel")
	f.BoolVar(
		&c.Force,
		"force",
		false,
		"force update checked out zettels, even if they will overwrite existing checkouts",
	)
}
