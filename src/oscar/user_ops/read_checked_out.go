package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/hinweisen"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	store_fs "github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type ReadCheckedOut struct {
	*umwelt.Umwelt
	store_fs.OptionsReadExternal
	AllowMissing bool
}

type ReadCheckedOutResults struct {
	Zettelen map[hinweis.Hinweis]zettel_checked_out.Zettel
}

func (op ReadCheckedOut) RunOneHinweis(
	h hinweis.Hinweis,
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
	possible store_fs.CwdFiles,
	w collections.WriterFunc[*zettel_checked_out.Zettel],
) (err error) {
	for _, p := range possible.Zettelen {
		var checked_out zettel_checked_out.Zettel

		var readFunc func() (zettel_checked_out.Zettel, error)

		switch {
		case p.Akte.Path == "":
			readFunc = func() (zettel_checked_out.Zettel, error) {
				return op.StoreWorkingDirectory().Read(p.Zettel.Path)
			}

		case p.Zettel.Path == "":
			readFunc = func() (zettel_checked_out.Zettel, error) {
				return op.StoreWorkingDirectory().ReadExternalZettelFromAktePath(p.Akte.Path)
			}

		default:
			//TODO validate that the zettel file points to the akte in the metadatei
			readFunc = func() (zettel_checked_out.Zettel, error) {
				return op.StoreWorkingDirectory().Read(p.Zettel.Path)
			}
		}

		if checked_out, err = readFunc(); err != nil {
			if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
				errors.Print("external zettel does not exist: %s", p)
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}

		}

		if err = w(&checked_out); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
