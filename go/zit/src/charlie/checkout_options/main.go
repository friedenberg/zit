package checkout_options

import (
	"flag"

	"code.linenisgreat.com/zit/src/bravo/checkout_mode"
)

type Options struct {
	Force           bool
	Path            Path
	ForceInlineAkte bool
	CheckoutMode    checkout_mode.Mode
	AllowConflicted bool
}

func (c *Options) AddToFlagSet(f *flag.FlagSet) {
	f.Var(&c.CheckoutMode, "mode", "mode for checking out the zettel")
	f.BoolVar(
		&c.Force,
		"force",
		false,
		"force update checked out zettels, even if they will overwrite existing checkouts",
	)
}
