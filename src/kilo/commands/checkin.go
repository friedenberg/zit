package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/golf/zettel_formats"
	checkout_store "github.com/friedenberg/zit/src/hotel/store_checkout"
	"github.com/friedenberg/zit/src/india/store_with_lock"
	"github.com/friedenberg/zit/src/juliett/user_ops"
)

type Checkin struct {
	Delete     bool
	IgnoreAkte bool
	All        bool
}

func init() {
	registerCommand(
		"checkin",
		func(f *flag.FlagSet) Command {
			c := &Checkin{}

			f.BoolVar(&c.Delete, "delete", false, "the checked-out file")
			f.BoolVar(&c.IgnoreAkte, "ignore-akte", false, "do not change the akte")
			f.BoolVar(&c.All, "all", false, "")

			return commandWithLockedStore{c}
		},
	)
}

func (c Checkin) RunWithLockedStore(
	s store_with_lock.Store,
	args ...string,
) (err error) {
	if c.All {
		if len(args) > 0 {
			stdprinter.Errf("Ignoring args because -all is set\n")
		}

		var possible checkout_store.CwdFiles

		if possible, err = user_ops.NewGetPossibleZettels(s.Umwelt).Run(s); err != nil {
			err = errors.Error(err)
			return
		}

		args = possible.Zettelen
	}

	checkinOptions := checkout_store.CheckinOptions{
		IncludeAkte: !c.IgnoreAkte,
		Format:      zettel_formats.Text{},
	}

	var readResults user_ops.ReadCheckedOutResults

	readOp := user_ops.ReadCheckedOut{
		Umwelt:  s.Umwelt,
		Options: checkinOptions,
	}

	if readResults, err = readOp.RunManyStrings(s, args...); err != nil {
		err = errors.Error(err)
		return
	}

	checkinOp := user_ops.Checkin{
		Umwelt:  s.Umwelt,
		Options: checkinOptions,
	}

	zettels := make([]stored_zettel.External, 0, len(readResults.Zettelen))

	for _, z := range readResults.Zettelen {
		zettels = append(zettels, z.External)
	}

	if _, err = checkinOp.Run(s, zettels...); err != nil {
		err = errors.Error(err)
		return
	}

	if c.Delete {
		deleteOp := user_ops.DeleteCheckout{
			Umwelt: s.Umwelt,
		}

		external := make(map[hinweis.Hinweis]stored_zettel.External)

		for h, z := range readResults.Zettelen {
			external[h] = z.External
		}

		if err = deleteOp.Run(s, external); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}
