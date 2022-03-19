package commands

import (
	"flag"

	"github.com/friedenberg/zit/foxtrot/stored_zettel"
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

func (c Checkin) Run(u _Umwelt, args ...string) (err error) {
	if c.All {
		if len(args) > 0 {
			_Errf("Ignoring args because -all is set\n")
		}

		getPossibleOp := user_ops.GetPossibleZettels{
			Umwelt: u,
		}

		var getPossibleResults user_ops.GetPossibleZettelsResults

		if getPossibleResults, err = getPossibleOp.Run(); err != nil {
			err = _Error(err)
			return
		}

		args = getPossibleResults.Hinweisen
	}

	checkinOptions := _ZettelsCheckinOptions{
		IncludeAkte: !c.IgnoreAkte,
		Format:      _ZettelFormatsText{},
	}

	var readResults user_ops.ReadCheckedOutResults

	readOp := user_ops.ReadCheckedOut{
		Umwelt:  u,
		Options: checkinOptions,
	}

	if readResults, err = readOp.Run(args...); err != nil {
		err = _Error(err)
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
		err = _Error(err)
		return
	}

	if c.Delete {
		deleteOp := user_ops.DeleteCheckout{}

		if err = deleteOp.Run(readResults.Zettelen); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}
