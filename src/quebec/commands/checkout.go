package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Checkout struct {
	CheckoutMode objekte.CheckoutMode
	Force        bool
}

func init() {
	registerCommandWithQuery(
		"checkout",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Checkout{}

			f.Var(&c.CheckoutMode, "mode", "mode for checking out the zettel")
			f.BoolVar(
				&c.Force,
				"force",
				false,
				"force update checked out zettels, even if they will overwrite existing checkouts",
			)

			return c
		},
	)
}

func (c Checkout) RunWithQuery(u *umwelt.Umwelt, ms kennung.MetaSet) (err error) {
	options := store_fs.CheckoutOptions{
		Force:        c.Force,
		CheckoutMode: c.CheckoutMode,
	}

	if err = u.StoreWorkingDirectory().CheckoutQuery(
		options,
		ms,
		func(co objekte.CheckedOutLike) (err error) {
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
