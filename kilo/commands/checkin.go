package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
	"github.com/friedenberg/zit/golf/checkout_store"
	"github.com/friedenberg/zit/hotel/zettels"
	"github.com/friedenberg/zit/juliett/user_ops"
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

			return c
		},
	)
}

func (c Checkin) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if c.All {
		if len(args) > 0 {
			stdprinter.Errf("Ignoring args because -all is set\n")
		}

		var possible checkout_store.CwdFiles

		if possible, err = user_ops.NewGetPossibleZettels(u).Run(); err != nil {
			err = errors.Error(err)
			return
		}

		args = possible.Zettelen
	}

	checkinOptions := zettels.CheckinOptions{
		IncludeAkte: !c.IgnoreAkte,
		Format:      zettel_formats.Text{},
	}

	var readResults user_ops.ReadCheckedOutResults

	readOp := user_ops.ReadCheckedOut{
		Umwelt:  u,
		Options: checkinOptions,
	}

	if readResults, err = readOp.Run(args...); err != nil {
		err = errors.Error(err)
		return
	}

	checkinOp := user_ops.Checkin{
		Umwelt:  u,
		Options: checkinOptions,
	}

	zettels := make([]stored_zettel.External, 0, len(readResults.Zettelen))

	for _, z := range readResults.Zettelen {
		zettels = append(zettels, z.External)
	}

	if _, err = checkinOp.Run(zettels...); err != nil {
		err = errors.Error(err)
		return
	}

	if c.Delete {
		deleteOp := user_ops.DeleteCheckout{
			Umwelt: u,
		}

		external := make(map[hinweis.Hinweis]stored_zettel.External)

		for h, z := range readResults.Zettelen {
			external[h] = z.External
		}

		if err = deleteOp.Run(external); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}
