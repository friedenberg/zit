package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/golf/zettel_formats"
	checkout_store "github.com/friedenberg/zit/src/hotel/store_checkout"
	"github.com/friedenberg/zit/src/india/store_with_lock"
	"github.com/friedenberg/zit/src/juliett/user_ops"
)

type Checkout struct {
	checkout_store.CheckoutMode
	Force bool
}

func init() {
	registerCommand(
		"checkout",
		func(f *flag.FlagSet) Command {
			c := &Checkout{}

			f.Var(&c.CheckoutMode, "mode", "mode for checking out the zettel")
			f.BoolVar(&c.Force, "force", false, "force update checked out zettels, even if they will overwrite existing checkouts")

			return commandWithLockedStore{commandWithHinweisen{c}}
		},
	)
}

func (c Checkout) RunWithHinweisen(s store_with_lock.Store, hins ...hinweis.Hinweis) (err error) {
	var readResults user_ops.ReadCheckedOutResults

	readOp := user_ops.ReadCheckedOut{
		Umwelt: s.Umwelt,
		OptionsReadExternal: checkout_store.OptionsReadExternal{
			Format: zettel_formats.Text{},
		},
	}

	if readResults, err = readOp.RunManyHinweisen(s, hins...); err != nil {
		logz.Print(err)
		err = errors.Error(err)
		return
	}

	toCheckOut := make([]hinweis.Hinweis, 0, len(hins))

	for h, cz := range readResults.Zettelen {
		if cz.External.Path == "" {
			toCheckOut = append(toCheckOut, h)
			continue
		}

		if cz.Internal.Zettel.Equals(cz.External.Zettel) {
			logz.Print(cz.Internal.Zettel)
			stdprinter.Outf("%s (already checked out)\n", cz.Internal.Named)
			continue
		}

		if c.Force || cz.External.Sha.IsNull() {
			toCheckOut = append(toCheckOut, h)
		} else {
			stdprinter.Errf("[%s] (external has changes)\n", h)
			continue
		}
	}

	checkoutOp := user_ops.Checkout{
		Umwelt: s.Umwelt,
		CheckoutOptions: checkout_store.CheckoutOptions{
			CheckoutMode: c.CheckoutMode,
			Format:       zettel_formats.Text{},
		},
	}

	if _, err = checkoutOp.RunManyHinweisen(s, toCheckOut...); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
