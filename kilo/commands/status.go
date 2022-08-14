package commands

import (
	"flag"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/stdprinter"
	"github.com/friedenberg/zit/golf/zettel_formats"
	checkout_store "github.com/friedenberg/zit/golf/store_checkout"
	"github.com/friedenberg/zit/india/store_with_lock"
	"github.com/friedenberg/zit/juliett/user_ops"
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

	var possible checkout_store.CwdFiles

	if possible, err = user_ops.NewGetPossibleZettels(s.Umwelt).Run(s); err != nil {
		err = errors.Error(err)
		return
	}

	args = possible.Zettelen

	options := checkout_store.CheckinOptions{
		IncludeAkte: true,
		Format:      zettel_formats.Text{},
	}

	var readResults user_ops.ReadCheckedOutResults

	readOp := user_ops.ReadCheckedOut{
		Umwelt:  s.Umwelt,
		Options: options,
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
