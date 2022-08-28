package user_ops

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/golf/hinweisen"
	store_working_directory "github.com/friedenberg/zit/src/hotel/store_working_directory"
	"github.com/friedenberg/zit/src/india/store_with_lock"
	"github.com/friedenberg/zit/src/india/zettel_checked_out"
)

type ReadCheckedOut struct {
	*umwelt.Umwelt
	store_working_directory.OptionsReadExternal
	AllowMissing bool
}

type ReadCheckedOutResults struct {
	Zettelen map[hinweis.Hinweis]zettel_checked_out.CheckedOut
}

func (op ReadCheckedOut) RunOneHinweis(
	s store_with_lock.Store,
	h hinweis.Hinweis,
) (zettel zettel_checked_out.CheckedOut, err error) {
	return op.RunOneString(s, h.String())
}

func (op ReadCheckedOut) RunOneString(
	s store_with_lock.Store,
	path string,
) (zettel zettel_checked_out.CheckedOut, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(op.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	if zettel, err = op.runOne(store, path); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (op ReadCheckedOut) RunMany(
	s store_with_lock.Store,
	possible store_working_directory.CwdFiles,
) (results []zettel_checked_out.CheckedOut, err error) {
	results = make([]zettel_checked_out.CheckedOut, 0, possible.Len())

	for _, p := range possible.Zettelen {
		var checked_out zettel_checked_out.CheckedOut

		if checked_out, err = op.runOne(s, p); err != nil {
			if errors.Is(err, hinweisen.ErrDoesNotExist) {
				//TODO log
				err = nil
			} else {
				err = errors.Error(err)
				return
			}

		}

		results = append(results, checked_out)
	}

	for _, p := range possible.Akten {
		var checked_out zettel_checked_out.CheckedOut

		if checked_out, err = s.StoreWorkingDirectory().ReadExternalZettelFromAktePath(p); err != nil {
			if errors.Is(err, hinweisen.ErrDoesNotExist) {
				//TODO log
				err = nil
			} else {
				err = errors.Error(err)
				return
			}
		}

		results = append(results, checked_out)
	}

	return
}

func (op ReadCheckedOut) runOne(
	store store_with_lock.Store,
	p string,
) (zettel zettel_checked_out.CheckedOut, err error) {
	return store.StoreWorkingDirectory().Read(p)
}
