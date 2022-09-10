package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_external"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/india/store_working_directory"
	"github.com/friedenberg/zit/src/juliett/umwelt"
	"github.com/friedenberg/zit/src/kilo/user_ops"
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

func (c Checkin) Run(
	s *umwelt.Umwelt,
	args ...string,
) (err error) {
	var pz store_working_directory.CwdFiles

	if c.All {
		if len(args) > 0 {
			errors.PrintErrf("Ignoring args because -all is set")
		}

		if pz, err = user_ops.NewGetPossibleZettels(s).Run(); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {

		pz.Zettelen = args
	}

	var readResults []zettel_checked_out.Zettel

	readOp := user_ops.ReadCheckedOut{
		Umwelt: s,
		OptionsReadExternal: store_working_directory.OptionsReadExternal{
			Format: zettel.Text{},
		},
	}

	if readResults, err = readOp.RunMany(pz); err != nil {
		err = errors.Wrap(err)
		return
	}

	checkinOp := user_ops.Checkin{
		Umwelt:              s,
		OptionsReadExternal: readOp.OptionsReadExternal,
	}

	zettels := make([]zettel_external.Zettel, 0, len(readResults))

	for _, z := range readResults {
		zettels = append(zettels, z.External)
	}

	if _, err = checkinOp.Run(zettels...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.Delete {
		deleteOp := user_ops.DeleteCheckout{
			Umwelt: s,
		}

		external := make([]zettel_external.Zettel, 0, len(readResults))

		for _, z := range readResults {
			external = append(external, z.External)
		}

		if err = deleteOp.Run(external); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
