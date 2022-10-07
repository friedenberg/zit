package commands

import (
	"flag"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/india/store_objekten"
	"github.com/friedenberg/zit/src/juliett/store_working_directory"
	"github.com/friedenberg/zit/src/kilo/umwelt"
	"github.com/friedenberg/zit/src/lima/user_ops"
)

type Status struct {
}

func init() {
	registerCommand(
		"status",
		func(f *flag.FlagSet) Command {
			c := &Status{}

			return c
		},
	)
}

func (c Status) Run(s *umwelt.Umwelt, args ...string) (err error) {
	var possible store_working_directory.CwdFiles

	switch {
	case len(args) > 0:
		errors.PrintErrf("Ignoring args")
		fallthrough

	default:
		if possible, err = store_working_directory.MakeCwdFiles(s.Standort().Cwd(), args...); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	options := store_working_directory.OptionsReadExternal{
		Format: zettel.Text{},
	}

	var readResultsSet zettel_checked_out.Set

	readOp := user_ops.ReadCheckedOut{
		Umwelt:              s,
		OptionsReadExternal: options,
	}

	if readResultsSet, err = readOp.RunMany(possible); err != nil {
		err = errors.Wrap(err)
		return
	}

	readResults := readResultsSet.ToSlice()

	sort.Slice(
		readResults,
		func(i, j int) bool {
			return readResults[i].External.ZettelFD.Path < readResults[j].External.ZettelFD.Path
		},
	)

	zp := s.PrinterOut()

	if !zp.IsEmpty() {
		err = zp.Error()
		return
	}

	for _, z := range readResults {
		zp.ZettelCheckedOut(z).Print()

		switch {
		case zp.IsEPIPE():
			zp.ClearErr()
			return

		case !zp.IsEmpty():
			err = zp.Error()
			return
		}
	}

	for _, ua := range possible.UnsureAkten {
		err = s.StoreObjekten().AkteExists(ua.Sha)

		switch {
		case err == nil:
			fallthrough
		case errors.Is(err, store_objekten.ErrNotFound{}):
			zp.FileUnrecognized(ua).Print()

			switch {
			case zp.IsEPIPE():
				zp.ClearErr()
				return

			case !zp.IsEmpty():
				err = zp.Error()
				return
			}

		case errors.Is(err, store_objekten.ErrAkteExists{}):
			err1 := err.(store_objekten.ErrAkteExists)
			zp.FileRecognized(ua, err1.Set).Print()
			err = zp.Error()

		default:
			err = errors.Wrapf(err, "%s", ua)
			return
		}
	}

	return
}
