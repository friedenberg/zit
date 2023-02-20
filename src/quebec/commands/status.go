package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/lima/store_objekten"
	"github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Status struct{}

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
	var possible cwd.CwdFiles

	switch {
	case len(args) > 0:
		errors.PrintErrf("Ignoring args")
		fallthrough

	default:
		if possible, err = cwd.MakeCwdFilesAll(
			s.Konfig(),
			s.Standort().Cwd(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	pcol := s.PrinterCheckedOutLike()

	if err = s.StoreWorkingDirectory().ReadFiles(
		possible,
		func(co objekte.CheckedOutLike) (err error) {
			if err = pcol(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	v := "Akten"

	if err = s.PrinterHeader()(&v); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, ua := range possible.UnsureAkten {
		err = s.StoreObjekten().AkteExists(ua.Sha)

		switch {
		case err == nil:
			fallthrough

		case errors.Is(err, objekte_store.ErrNotFound{}):
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
