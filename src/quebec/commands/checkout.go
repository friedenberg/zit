package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Checkout struct {
	store_fs.CheckoutMode
	Or    bool
	Force bool
}

func init() {
	registerCommand(
		"checkout",
		func(f *flag.FlagSet) Command {
			c := &Checkout{}

			f.BoolVar(&c.Or, "or", false, "allow optional criteria instead of required")
			f.Var(&c.CheckoutMode, "mode", "mode for checking out the zettel")
			f.BoolVar(&c.Force, "force", false, "force update checked out zettels, even if they will overwrite existing checkouts")

			return commandWithIds{
				CommandWithIds: c,
			}
		},
	)
}

func (c Checkout) RunWithIds(u *umwelt.Umwelt, ids kennung.Set) (err error) {
	options := store_fs.CheckoutOptions{
		CheckoutMode: c.CheckoutMode,
	}

	query := zettel.WriterIds{
		Filter: kennung.Filter{
			Set: ids,
			Or:  c.Or,
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
