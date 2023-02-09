package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Checkout struct {
	store_fs.CheckoutMode
	Force bool
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
		CheckoutMode: c.CheckoutMode,
	}

	ids, ok := ms.Get(gattung.Zettel)

	if !ok {
		return
	}

	query := zettel.WriterIds{
		Filter: kennung.Filter{
			Set: ids,
		},
	}

	if _, err = u.StoreWorkingDirectory().Checkout(
		options,
		query.WriteZettelTransacted,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
