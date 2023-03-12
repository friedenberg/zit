package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/lima/store_objekten"
	"github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Status struct{}

func init() {
	registerCommandWithCwdQuery(
		"status",
		func(f *flag.FlagSet) CommandWithCwdQuery {
			c := &Status{}

			return c
		},
	)
}

func (c Status) RunWithCwdQuery(
	u *umwelt.Umwelt,
	ms kennung.MetaSet,
	possible cwd.CwdFiles,
) (err error) {
	pcol := u.PrinterCheckedOutLike()

	if err = u.StoreWorkingDirectory().ReadFiles(
		possible,
		ms,
		iter.MakeChain(
			objekte.MakeFilterFromMetaSet(ms),
			func(co objekte.CheckedOutLike) (err error) {
				if err = pcol(co); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if len(possible.UnsureAkten) != 0 {

		v := "Akten"

		if err = u.PrinterHeader()(&v); err != nil {
			err = errors.Wrap(err)
			return
		}

		for _, ua := range possible.UnsureAkten {
			err = u.StoreObjekten().AkteExists(ua.Sha)

			switch {
			case err == nil:
				fallthrough

			case errors.Is(err, objekte_store.ErrNotFound{}):
				err = u.PrinterFileNotRecognized()(&ua)

			case errors.Is(err, store_objekten.ErrAkteExists{}):
				err1 := err.(store_objekten.ErrAkteExists)
				fr := store_fs.FileRecognized{
					FD:         ua,
					Recognized: err1.MutableSet,
				}

				err = u.PrinterFileRecognized()(&fr)

			default:
				err = errors.Wrapf(err, "%s", ua)
				return
			}
		}
	}

	return
}
