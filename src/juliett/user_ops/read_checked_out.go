package user_ops

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/charlie/sha"
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

func (op ReadCheckedOut) RunMany(
	s store_with_lock.Store,
	possible store_checkout.CwdFiles,
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

		if checked_out, err = op.runOneAkte(s, p); err != nil {
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
			err = nil
		} else {
			err = errors.Error(err)
			return
		}
	}

	zettel.DetermineState()

	if zettel.State > zettel_checked_out.StateExistsAndSame {
		exSha := zettel.External.Stored.Sha
		zettel.Matches.Zettelen, _ = store.Zettels().ReadZettelSha(exSha)

		exAkteSha := zettel.External.Stored.Zettel.Akte
		zettel.Matches.Akten, _ = store.Zettels().ReadAkteSha(exAkteSha)

		bez := zettel.External.Stored.Zettel.Bezeichnung.String()
		zettel.Matches.Bezeichnungen, _ = store.Zettels().ReadBezeichnung(bez)
	}

	return
}

func (op ReadCheckedOut) runOneAkte(
	store store_with_lock.Store,
	p string,
) (zettel zettel_checked_out.CheckedOut, err error) {
	stdprinter.Out(p)
	var akteSha sha.Sha

	if akteSha, err = store.CheckoutStore().AkteShaFromPath(p); err != nil {
		err = errors.Error(err)
		return
	}

	zettel.External.AktePath = p
	zettel.External.Named.Stored.Zettel.Akte = akteSha
	zettel.State = zettel_checked_out.StateAkte
	zettel.Matches.Akten, _ = store.Zettels().ReadAkteSha(akteSha)

	return
}
