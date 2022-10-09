package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_external"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/kilo/store_working_directory"
	"github.com/friedenberg/zit/src/mike/umwelt"
	"github.com/friedenberg/zit/src/november/user_ops"
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

	switch {
	case c.All && len(args) > 0:
		errors.PrintErrf("Ignoring args because -all is set")
		fallthrough

	default:
		if pz, err = store_working_directory.MakeCwdFiles(s.Standort().Cwd(), args...); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	var readResults zettel_checked_out.Set

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

	zettels := make([]zettel_external.Zettel, 0, readResults.Len())

	err = readResults.Each(
		func(zco zettel_checked_out.Zettel) (err error) {
			zettels = append(zettels, zco.External)
			return
		},
	)

	if _, err = checkinOp.Run(zettels...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.Delete {
		deleteOp := user_ops.DeleteCheckout{
			Umwelt: s,
		}

		external := make([]zettel_external.Zettel, 0, readResults.Len())

		err = readResults.Each(
			func(zco zettel_checked_out.Zettel) (err error) {
				external = append(external, zco.External)
				return
			},
		)

		if err = deleteOp.Run(external); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
