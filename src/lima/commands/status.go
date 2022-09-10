package commands

import (
	"flag"
	"os"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/hotel/store_objekten"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/india/store_working_directory"
	"github.com/friedenberg/zit/src/juliett/umwelt"
	"github.com/friedenberg/zit/src/juliett/zettel_printer"
	"github.com/friedenberg/zit/src/kilo/user_ops"
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
	if len(args) > 0 {
		errors.PrintErrf("args provided will be ignored")
	}

	var possible store_working_directory.CwdFiles

	if possible, err = user_ops.NewGetPossibleZettels(s).Run(); err != nil {
		err = errors.Wrap(err)
		return
	}

	options := store_working_directory.OptionsReadExternal{
		Format: zettel.Text{},
	}

	var readResults []zettel_checked_out.Zettel

	readOp := user_ops.ReadCheckedOut{
		Umwelt:              s,
		OptionsReadExternal: options,
	}

	if readResults, err = readOp.RunMany(possible); err != nil {
		err = errors.Wrap(err)
		return
	}

	sort.Slice(
		readResults,
		func(i, j int) bool {
			return readResults[i].External.ZettelFD.Path < readResults[j].External.ZettelFD.Path
		},
	)

	zp := zettel_printer.Make(s.StoreObjekten(), os.Stdout)
	zp.ShouldAbbreviateHinweisen = true

	if !zp.IsEmpty() {
		err = zp.Error()
		return
	}

	for _, z := range readResults {
		//TODO move to printer
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
