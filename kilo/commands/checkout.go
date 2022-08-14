package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/stdprinter"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
	checkout_store "github.com/friedenberg/zit/golf/store_checkout"
	"github.com/friedenberg/zit/india/store_with_lock"
	"github.com/friedenberg/zit/juliett/user_ops"
)

type Checkout struct {
	IncludeAkte bool
	Force       bool
}

func init() {
	registerCommand(
		"checkout",
		func(f *flag.FlagSet) Command {
			c := &Checkout{}

			f.BoolVar(&c.IncludeAkte, "include-akte", false, "check out akte as well")
			f.BoolVar(&c.Force, "force", false, "force update checked out zettels, even if they will overwrite existing checkouts")

			return commandWithLockedStore{commandWithHinweisen{c}}
		},
	)
}

func (c Checkout) RunWithHinweisen(s store_with_lock.Store, hins ...hinweis.Hinweis) (err error) {
	// getHinweisenOp := user_ops.GetAllHinweisen{
	// 	Umwelt: u,
	// }

	// var getHinweisenResults user_ops.GetAllHinweisenResults

	// if getHinweisenResults, err = getHinweisenOp.Run(); err != nil {
	// 	err = errors.Error(err)
	// 	return
	// }

	// hins = getHinweisenResults.Hinweisen

	checkinOptions := checkout_store.CheckinOptions{
		IgnoreMissingHinweis: true,
		AddMdExtension:       true,
		IncludeAkte:          c.IncludeAkte,
		Format:               zettel_formats.Text{},
	}

	var readResults user_ops.ReadCheckedOutResults

	readOp := user_ops.ReadCheckedOut{
		Umwelt:  s.Umwelt,
		Options: checkinOptions,
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

	options := checkout_store.CheckinOptions{
		IncludeAkte: c.IncludeAkte,
		Format:      zettel_formats.Text{},
	}

	checkoutOp := user_ops.Checkout{
		Umwelt:  s.Umwelt,
		Options: options,
	}

	if _, err = checkoutOp.RunManyHinweisen(s, toCheckOut...); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
