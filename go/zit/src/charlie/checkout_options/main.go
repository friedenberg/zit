package checkout_options

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
)

type Options struct {
	CheckoutMode checkout_mode.Mode
	OptionsWithoutMode
}

type OptionsWithoutMode struct {
	Force                bool
	AllowConflicted      bool
	StoreSpecificOptions any
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
