package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/kilo/store_working_directory"
	"github.com/friedenberg/zit/src/lima/store_with_lock"
	"github.com/friedenberg/zit/src/mike/user_ops"
)

type Checkout struct {
	store_working_directory.CheckoutMode
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
	var readResults []zettel_checked_out.Zettel

	readOp := user_ops.ReadCheckedOut{
		Umwelt: s.Umwelt,
		OptionsReadExternal: store_working_directory.OptionsReadExternal{
			Format: zettel.Text{},
		},
	}

	var pz store_working_directory.CwdFiles

	for _, h := range hins {
		pz.Zettelen = append(pz.Zettelen, h.String())
	}

	if readResults, err = readOp.RunMany(s, pz); err != nil {
		errors.Print(err)
		err = errors.Wrap(err)
		return
	}

	toCheckOut := make([]hinweis.Hinweis, 0, len(hins))

	for _, cz := range readResults {
		if cz.External.ZettelFD.Path == "" {
			toCheckOut = append(toCheckOut, cz.Internal.Named.Hinweis)
			continue
		}

		if cz.Internal.Named.Stored.Zettel.Equals(cz.External.Named.Stored.Zettel) {
			errors.Print(cz.Internal.Named.Stored.Zettel)
			errors.PrintOutf("%s (already checked out)", cz.Internal.Named)
			continue
		}

		if c.Force || cz.State == zettel_checked_out.StateEmpty {
			toCheckOut = append(toCheckOut, cz.Internal.Named.Hinweis)
		} else if cz.State == zettel_checked_out.StateExistsAndSame {
			errors.PrintOutf("%s (already checked out)", cz.Internal.Named)
			continue
		} else if cz.State == zettel_checked_out.StateExistsAndDifferent {
			errors.PrintOutf("%s (external has changes)", cz.Internal.Named)
			continue
		} else {
			errors.PrintOutf("%s (unknown state)", cz.Internal.Named)
			continue
		}
	}

	checkoutOp := user_ops.Checkout{
		Umwelt: s.Umwelt,
		CheckoutOptions: store_working_directory.CheckoutOptions{
			CheckoutMode: c.CheckoutMode,
			Format:       zettel.Text{},
		},
	}

	if _, err = checkoutOp.RunManyHinweisen(s, toCheckOut...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
