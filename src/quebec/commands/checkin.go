package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/kilo/cwd_files"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/zettel_checked_out"
	"github.com/friedenberg/zit/src/november/store_fs"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/user_ops"
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
	var pz cwd_files.CwdFiles

	switch {
	case c.All && len(args) > 0:
		errors.PrintErrf("Ignoring args because -all is set")
		fallthrough

	case c.All:
		if pz, err = cwd_files.MakeCwdFilesAll(s.Konfig(), s.Standort().Cwd()); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		if pz, err = cwd_files.MakeCwdFilesExactly(
			s.Konfig(),
			s.Standort().Cwd(),
			args...,
		); err != nil {
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
