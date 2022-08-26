package commands

import (
	"flag"
	"sort"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/golf/zettel_formats"
	store_working_directory "github.com/friedenberg/zit/src/hotel/store_working_directory"
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

	var possible store_working_directory.CwdFiles

	if possible, err = user_ops.NewGetPossibleZettels(s.Umwelt).Run(s); err != nil {
		err = errors.Error(err)
		return
	}

	options := store_working_directory.OptionsReadExternal{
		Format: zettel_formats.Text{},
	}

	var readResults []zettel_checked_out.CheckedOut

	readOp := user_ops.ReadCheckedOut{
		Umwelt:              s.Umwelt,
		OptionsReadExternal: options,
	}

	if readResults, err = readOp.RunMany(s, possible); err != nil {
		err = errors.Error(err)
		return
	}

	sort.Slice(
		readResults,
		func(i, j int) bool {
			return readResults[i].External.Path < readResults[j].External.Path
		},
	)

	for _, z := range readResults {
		if z.State == zettel_checked_out.StateEmpty {
			stdprinter.Outf("%#v\n", z.External)
		} else {
			stdprinter.Out(z)
		}
	}

	return
}
