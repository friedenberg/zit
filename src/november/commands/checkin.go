package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/kilo/store_working_directory"
	"github.com/friedenberg/zit/src/lima/store_with_lock"
	"github.com/friedenberg/zit/src/mike/user_ops"
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
	var pz store_working_directory.CwdFiles

	if c.All {
		if len(args) > 0 {
			stdprinter.Errf("Ignoring args because -all is set\n")
		}

		if pz, err = user_ops.NewGetPossibleZettels(s.Umwelt).Run(s); err != nil {
			err = errors.Error(err)
			return
		}
	} else {

		pz.Zettelen = args
	}

	var readResults []zettel_checked_out.Zettel

	readOp := user_ops.ReadCheckedOut{
		Umwelt: s.Umwelt,
		OptionsReadExternal: store_working_directory.OptionsReadExternal{
			Format: zettel.Text{},
		},
	}

	if readResults, err = readOp.RunMany(s, pz); err != nil {
		err = errors.Error(err)
		return
	}

	checkinOp := user_ops.Checkin{
		Umwelt:              s.Umwelt,
		OptionsReadExternal: readOp.OptionsReadExternal,
	}

	zettels := make([]zettel_external.Zettel, 0, len(readResults))

	for _, z := range readResults {
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

		external := make([]zettel_external.Zettel, 0, len(readResults))

		for _, z := range readResults {
			external = append(external, z.External)
		}

		if err = deleteOp.Run(s, external); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}
