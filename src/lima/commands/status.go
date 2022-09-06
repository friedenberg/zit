package commands

import (
	"flag"
	"sort"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/store_objekten"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/india/store_working_directory"
	"github.com/friedenberg/zit/src/juliett/store_with_lock"
	"github.com/friedenberg/zit/src/kilo/user_ops"
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
		errors.PrintErrf("args provided will be ignored")
	}

	var possible store_working_directory.CwdFiles

	if possible, err = user_ops.NewGetPossibleZettels(s.Umwelt).Run(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	options := store_working_directory.OptionsReadExternal{
		Format: zettel.Text{},
	}

	var readResults []zettel_checked_out.Zettel

	readOp := user_ops.ReadCheckedOut{
		Umwelt:              s.Umwelt,
		OptionsReadExternal: options,
	}

	if readResults, err = readOp.RunMany(s, possible); err != nil {
		err = errors.Wrap(err)
		return
	}

	sort.Slice(
		readResults,
		func(i, j int) bool {
			return readResults[i].External.ZettelFD.Path < readResults[j].External.ZettelFD.Path
		},
	)

	for _, z := range readResults {
		if z.State == zettel_checked_out.StateEmpty {
			errors.PrintOutf("%#v", z.External)
		} else {
			errors.PrintOut(z)
		}
	}

	for _, ua := range possible.UnsureAkten {
		if err = s.StoreObjekten().AkteExists(ua.Sha); err == nil {
			if err = errors.PrintOutf("%s (not recognized)", ua); err != nil {
				err = errors.IsAsNilOrWrapf(
					err,
					syscall.EPIPE,
					"%s",
					ua,
				)
			}
		} else if errors.Is(err, store_objekten.ErrNotFound{}) {
			if err = errors.PrintOutf("%s (not recognized)", ua); err != nil {
				err = errors.IsAsNilOrWrapf(
					err,
					syscall.EPIPE,
					"%s",
					ua,
				)
			}
		} else if errors.Is(err, store_objekten.ErrAkteExists{}) {
			err1 := err.(store_objekten.ErrAkteExists)
			errors.PrintOutf("%s (Akte recognized)", ua)
			err1.Set.Each(
				func(tz1 zettel_transacted.Zettel) (err error) {
					//TODO eliminate zettels marked as duplicates / hidden
					if err = errors.PrintOutf("\t%s", tz1.Named); err != nil {
						err = errors.IsAsNilOrWrapf(
							err,
							syscall.EPIPE,
							"%s",
							tz1.Named,
						)

						return
					}

					return
				},
			)

			err = nil

		} else {
			err = errors.Wrapf(err, "%s", ua)
			return
		}
	}

	return
}
