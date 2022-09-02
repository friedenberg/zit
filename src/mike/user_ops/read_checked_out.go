package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/golf/hinweisen"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/kilo/store_working_directory"
	"github.com/friedenberg/zit/src/lima/store_with_lock"
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
	s store_with_lock.Store,
	h hinweis.Hinweis,
) (zettel zettel_checked_out.Zettel, err error) {
	return op.RunOneString(s, h.String())
}

func (op ReadCheckedOut) RunOneString(
	s store_with_lock.Store,
	p string,
) (zettel zettel_checked_out.Zettel, err error) {
	if zettel, err = s.StoreWorkingDirectory().Read(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op ReadCheckedOut) RunMany(
	s store_with_lock.Store,
	possible store_working_directory.CwdFiles,
	//TODO switch to zettel_checked_out.Set
) (results []zettel_checked_out.Zettel, err error) {
	results = make([]zettel_checked_out.Zettel, 0, possible.Len())

	for _, p := range possible.Zettelen {
		var checked_out zettel_checked_out.Zettel

		if checked_out, err = s.StoreWorkingDirectory().Read(p); err != nil {
			if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
				//TODO log
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}

		}

		results = append(results, checked_out)
	}

	for _, p := range possible.Akten {
		var checked_out zettel_checked_out.Zettel

		if checked_out, err = s.StoreWorkingDirectory().ReadExternalZettelFromAktePath(p); err != nil {
			if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
				//TODO log
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		results = append(results, checked_out)
	}

	return
}
