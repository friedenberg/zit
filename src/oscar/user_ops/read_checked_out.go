package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/hinweisen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/cwd"
	"github.com/friedenberg/zit/src/lima/zettel_checked_out"
	"github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type ReadCheckedOut struct {
	*umwelt.Umwelt
	store_fs.OptionsReadExternal
	AllowMissing bool
}

type ReadCheckedOutResults struct {
	Zettelen map[kennung.Hinweis]zettel_checked_out.Zettel
}

func (op ReadCheckedOut) RunOneHinweis(
	h kennung.Hinweis,
) (zettel zettel_checked_out.Zettel, err error) {
	return op.RunOneString(h.String())
}

func (op ReadCheckedOut) RunOneString(
	p string,
) (zettel zettel_checked_out.Zettel, err error) {
	if zettel, err = op.StoreWorkingDirectory().Read(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op ReadCheckedOut) RunMany(
	possible cwd.CwdFiles,
	w schnittstellen.FuncIter[*zettel_checked_out.Zettel],
) (err error) {
	for _, p := range possible.Zettelen {
		var checked_out zettel_checked_out.Zettel

		var readFunc func() (zettel_checked_out.Zettel, error)

		switch {
		case p.AkteFD.Path == "":
			readFunc = func() (zettel_checked_out.Zettel, error) {
				return op.StoreWorkingDirectory().Read(p.FD.Path)
			}

		case p.FD.Path == "":
			readFunc = func() (zettel_checked_out.Zettel, error) {
				return op.StoreWorkingDirectory().ReadExternalZettelFromAktePath(
					p.AkteFD.Path,
				)
			}

		default:
			// TODO-P3 validate that the zettel file points to the akte in the metadatei
			readFunc = func() (zettel_checked_out.Zettel, error) {
				return op.StoreWorkingDirectory().Read(p.FD.Path)
			}
		}

		defer func() {
			if e := recover(); e != nil {
				errors.Err().Printf("Path: %s", p)
				panic(e)
			}
		}()

		if checked_out, err = readFunc(); err != nil {
			// TODO-P3 decide if error handling like this is ok
			if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
				errors.Err().Printf("external zettel does not exist: %s", p)
			} else {
				errors.Err().Print(err)
			}

			err = nil
			continue
		}

		if err = w(&checked_out); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
