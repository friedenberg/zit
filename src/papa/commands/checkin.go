package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	store_fs "github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/oscar/user_ops"
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
	var pz store_fs.CwdFiles

	switch {
	case c.All && len(args) > 0:
		errors.PrintErrf("Ignoring args because -all is set")
		fallthrough

	case c.All:
		if pz, err = store_fs.MakeCwdFilesAll(s.Konfig().Compiled, s.Standort().Cwd()); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		if pz, err = store_fs.MakeCwdFilesExactly(s.Konfig().Compiled, s.Standort().Cwd(), args...); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	readOp := user_ops.ReadCheckedOut{
		Umwelt: s,
		OptionsReadExternal: store_fs.OptionsReadExternal{
			Format: zettel.Text{},
		},
	}

	readResults := zettel_checked_out.MakeMutableSetUnique(0)

	if err = readOp.RunMany(pz, readResults.Add); err != nil {
		err = errors.Wrap(err)
		return
	}

	checkinOp := user_ops.Checkin{
		Umwelt:              s,
		OptionsReadExternal: readOp.OptionsReadExternal,
	}

	zettels := make([]zettel_external.Zettel, 0, readResults.Len())

	err = readResults.Each(
		func(zco *zettel_checked_out.Zettel) (err error) {
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

		external := zettel_external.MakeMutableSetUniqueFD()

		err = readResults.Each(
			func(zco *zettel_checked_out.Zettel) (err error) {
				return external.Add(&zco.External)
			},
		)

		if err = deleteOp.Run(external); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}