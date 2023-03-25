package store_fs

import (
	"flag"

	"github.com/friedenberg/zit/src/lima/store_objekten"
)

type CheckoutOptions store_objekten.CheckoutOptions

func (c *CheckoutOptions) AddToFlagSet(f *flag.FlagSet) {
	f.Var(&c.CheckoutMode, "mode", "mode for checking out the zettel")
	f.BoolVar(
		&c.Force,
		"force",
		false,
		"force update checked out zettels, even if they will overwrite existing checkouts",
	)
}
