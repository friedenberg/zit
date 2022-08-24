package commands

import (
	"flag"
	"sort"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/golf/zettel_formats"
	store_checkout "github.com/friedenberg/zit/src/hotel/store_checkout"
	"github.com/friedenberg/zit/src/india/store_with_lock"
	"github.com/friedenberg/zit/src/india/zettel_checked_out"
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

	sl := make([]zettel_checked_out.CheckedOut, 0, len(readResults.Zettelen))

	for _, z := range readResults.Zettelen {
		sl = append(sl, z)
	}

	sort.Slice(
		sl,
		func(i, j int) bool {
			return sl[i].External.Path < sl[j].External.Path
		},
	)

	for _, z := range sl {
		stdprinter.Out(z)
	}

	return
}
