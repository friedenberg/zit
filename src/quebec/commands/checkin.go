package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/cwd"
	"github.com/friedenberg/zit/src/kilo/zettel_external"
	"github.com/friedenberg/zit/src/lima/zettel_checked_out"
	"github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/oscar/user_ops"
)

type Checkin struct {
	Delete     bool
	IgnoreAkte bool
}

func init() {
	registerCommandWithQuery(
		"checkin",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Checkin{}

			f.BoolVar(&c.Delete, "delete", false, "the checked-out file")
			f.BoolVar(&c.IgnoreAkte, "ignore-akte", false, "do not change the akte")

			return c
		},
	)
}

func (c Checkin) RunWithQuery(
	s *umwelt.Umwelt,
	ms kennung.MetaSet,
) (err error) {
	var pz cwd.CwdFiles

	if pz, err = cwd.MakeCwdFilesMetaSet(
		s.Konfig(),
		s.Standort().Cwd(),
		ms,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	readOp := user_ops.ReadCheckedOut{
		Umwelt:              s,
		OptionsReadExternal: store_fs.OptionsReadExternal{},
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
