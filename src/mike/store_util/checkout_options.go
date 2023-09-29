package store_util

import (
	"flag"

	"github.com/friedenberg/zit/src/bravo/checkout_mode"
)

type CheckoutOptions struct {
	Force        bool
	CheckoutMode checkout_mode.Mode
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
