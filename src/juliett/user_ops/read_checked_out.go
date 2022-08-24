package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/golf/hinweisen"
	"github.com/friedenberg/zit/src/hotel/store_checkout"
	"github.com/friedenberg/zit/src/hotel/store_objekten"
	"github.com/friedenberg/zit/src/india/store_with_lock"
	"github.com/friedenberg/zit/src/india/zettel_checked_out"
)

type ReadCheckedOut struct {
	*umwelt.Umwelt
	store_checkout.OptionsReadExternal
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

func (op ReadCheckedOut) RunManyHinweisen(
	s store_with_lock.Store,
	hins ...hinweis.Hinweis,
) (results ReadCheckedOutResults, err error) {
	ss := make([]string, len(hins))

	for i, _ := range ss {
		ss[i] = hins[i].String()
	}

	return op.RunManyStrings(s, ss...)
}

func (op ReadCheckedOut) RunManyStrings(
	s store_with_lock.Store,
	paths ...string,
) (results ReadCheckedOutResults, err error) {
	results.Zettelen = make(map[hinweis.Hinweis]zettel_checked_out.CheckedOut)
	for _, p := range paths {
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

		logz.Print(checked_out.External.Hinweis)
		results.Zettelen[checked_out.External.Hinweis] = checked_out
	}

	return
}

func (op ReadCheckedOut) runOne(
	store store_with_lock.Store,
	p string,
) (zettel zettel_checked_out.CheckedOut, err error) {
	if zettel.External, err = store.CheckoutStore().Read(p); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Error(err)
			return
		}
	}

	if zettel.Internal, err = store.Zettels().Read(zettel.External.Hinweis); err != nil {
		if errors.Is(err, store_objekten.ErrNotFound{}) {
			exSha := zettel.External.Stored.Sha
			zettel.Matches, _ = store.Zettels().ReadZettelSha(exSha)
			err = nil
		} else {
			err = errors.Error(err)
			return
		}
	}

	zettel.DetermineState()

	return
}
