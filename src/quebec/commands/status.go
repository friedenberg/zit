package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/kilo/cwd_files"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/store_objekten"
	"github.com/friedenberg/zit/src/mike/zettel_checked_out"
	"github.com/friedenberg/zit/src/november/store_fs"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/user_ops"
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
	var possible cwd_files.CwdFiles

	switch {
	case len(args) > 0:
		errors.PrintErrf("Ignoring args")
		fallthrough

	default:
		if possible, err = cwd_files.MakeCwdFilesAll(
			s.Konfig(),
			s.Standort().Cwd(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	options := store_fs.OptionsReadExternal{
		Format: zettel.MakeTextParser(
			s.StoreObjekten(),
			nil,
		),
	}

	readOp := user_ops.ReadCheckedOut{
		Umwelt:              s,
		OptionsReadExternal: options,
	}

	readResultsSet := zettel_checked_out.MakeMutableSetUnique(0)

	if err = readOp.RunMany(possible, readResultsSet.Add); err != nil {
		err = errors.Wrap(err)
		return
	}

	v := "Typen"

	if err = s.PrinterHeader()(&v); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, p := range possible.Typen {
		if err = s.StoreWorkingDirectory().ReadTyp(p); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = s.PrinterTypCheckedOut("same")(p); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	v = "Zettelen"

	if err = s.PrinterHeader()(&v); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO-P4 use right mode
	if err = readResultsSet.Each(s.PrinterZettelCheckedOut(zettel_checked_out.ModeZettelAndAkte)); err != nil {
		err = errors.Wrap(err)
		return
	}

	v = "Akten"

	if err = s.PrinterHeader()(&v); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, ua := range possible.UnsureAkten {
		err = s.StoreObjekten().AkteExists(ua.Sha)

		switch {
		case err == nil:
			fallthrough

		case errors.Is(err, store_objekten.ErrNotFound{}):
			err = s.PrinterFileNotRecognized()(&ua)

		case errors.Is(err, store_objekten.ErrAkteExists{}):
			err1 := err.(store_objekten.ErrAkteExists)
			fr := store_fs.FileRecognized{
				FD:         ua,
				Recognized: err1.MutableSet,
			}

			err = s.PrinterFileRecognized()(&fr)

		default:
			err = errors.Wrapf(err, "%s", ua)
			return
		}
	}

	return
}
