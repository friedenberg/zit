package commands

import (
	"flag"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/lima/store_objekten"
	store_fs "github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/oscar/user_ops"
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
	var possible store_fs.CwdFiles

	switch {
	case len(args) > 0:
		errors.PrintErrf("Ignoring args")
		fallthrough

	default:
		if possible, err = store_fs.MakeCwdFilesAll(s.Konfig().Compiled, s.Standort().Cwd()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	options := store_fs.OptionsReadExternal{
		Format: zettel.Text{},
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

	readResults := readResultsSet.Elements()

	sort.Slice(
		readResults,
		func(i, j int) bool {
			return readResults[i].External.ZettelFD.Path < readResults[j].External.ZettelFD.Path
		},
	)

	if err = readResultsSet.Each(s.PrinterZettelCheckedOut()); err != nil {
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
				File:       ua,
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
