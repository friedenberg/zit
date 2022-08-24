package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/golf/zettel_formats"
	store_checkout "github.com/friedenberg/zit/src/hotel/store_checkout"
	"github.com/friedenberg/zit/src/india/store_with_lock"
	"github.com/friedenberg/zit/src/juliett/user_ops"
)

type Status struct {
}

func init() {
	registerCommand(
		"status",
		func(f *flag.FlagSet) Command {
			c := &Status{}

			return commandWithLockedStore{c}
		},
	)
}

func (c Status) RunWithLockedStore(s store_with_lock.Store, args ...string) (err error) {
	if len(args) > 0 {
		stdprinter.Errf("args provided will be ignored")
	}

	var possible store_checkout.CwdFiles

	if possible, err = user_ops.NewGetPossibleZettels(s.Umwelt).Run(s); err != nil {
		err = errors.Error(err)
		return
	}

	args = possible.Zettelen

	options := store_checkout.OptionsReadExternal{
		Format: zettel_formats.Text{},
	}

	var readResults user_ops.ReadCheckedOutResults

	readOp := user_ops.ReadCheckedOut{
		Umwelt:              s.Umwelt,
		OptionsReadExternal: options,
	}

	if readResults, err = readOp.RunManyStrings(s, args...); err != nil {
		err = errors.Error(err)
		return
	}

	for h, z := range readResults.Zettelen {
		if z.Internal.Zettel.Equals(z.External.Zettel) {
			continue
		}

		stdprinter.Outf("[%s] (different)\n", h)
	}

	return
}
