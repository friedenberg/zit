package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/delta/hinweisen"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/india/store_working_directory"
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

type ReadCheckedOut struct {
	*umwelt.Umwelt
	store_working_directory.OptionsReadExternal
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
	possible store_working_directory.CwdFiles,
	//TODO switch to zettel_checked_out.Set
) (results zettel_checked_out.Set, err error) {
  results = zettel_checked_out.MakeSetUnique(possible.Len())

	for _, p := range possible.Zettelen {
		var checked_out zettel_checked_out.Zettel

		if p.Zettel.Path == "" {
			continue
		}

		if checked_out, err = op.StoreWorkingDirectory().Read(p.Zettel.Path); err != nil {
			if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
				//TODO log
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}

		}

    results.Add(checked_out)
	}

	for _, p := range possible.Zettelen {
		var checked_out zettel_checked_out.Zettel

		if p.Akte.Path == "" {
			continue
		}

		if checked_out, err = op.StoreWorkingDirectory().ReadExternalZettelFromAktePath(p.Akte.Path); err != nil {
			if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
				//TODO log
				err = nil
			} else {
				err = errors.Wrapf(err, "akte path: %s", p.Akte.Path)
				return
			}
		}

    results.Add(checked_out)
	}

	return
}
