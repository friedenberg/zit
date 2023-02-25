package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/hinweisen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type ReadCheckedOut struct {
	*umwelt.Umwelt
	store_fs.OptionsReadExternal
	AllowMissing bool
}

type ReadCheckedOutResults struct {
	Zettelen map[kennung.Hinweis]zettel.CheckedOut
}

func (op ReadCheckedOut) RunOneHinweis(
	h kennung.Hinweis,
) (zettel zettel.CheckedOut, err error) {
	return op.RunOneString(h.String())
}

func (op ReadCheckedOut) RunOneString(
	p string,
) (zettel zettel.CheckedOut, err error) {
	if zettel, err = op.StoreWorkingDirectory().Read(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op ReadCheckedOut) RunMany(
	possible cwd.CwdFiles,
	w schnittstellen.FuncIter[*zettel.CheckedOut],
) (err error) {
	if err = possible.Zettelen.Each(
		func(p cwd.Zettel) (err error) {
			var checked_out zettel.CheckedOut

			var readFunc func() (zettel.CheckedOut, error)

			switch {
			case p.GetAkteFD().Path == "":
				readFunc = func() (zettel.CheckedOut, error) {
					return op.StoreWorkingDirectory().Read(p.GetObjekteFD().Path)
				}

			case p.GetObjekteFD().Path == "":
				readFunc = func() (zettel.CheckedOut, error) {
					return op.StoreWorkingDirectory().ReadExternalZettelFromAktePath(
						p.GetAkteFD().Path,
					)
				}

			default:
				// TODO-P3 validate that the zettel file points to the akte in the metadatei
				readFunc = func() (zettel.CheckedOut, error) {
					return op.StoreWorkingDirectory().Read(p.GetObjekteFD().Path)
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
				return
			}

			if err = w(&checked_out); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
