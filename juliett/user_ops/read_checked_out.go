package user_ops

import (
	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/delta/hinweis"
	"github.com/friedenberg/zit/echo/umwelt"
	"github.com/friedenberg/zit/golf/hinweisen"
	checkout_store "github.com/friedenberg/zit/golf/store_checkout"
	"github.com/friedenberg/zit/golf/stored_zettel"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type ReadCheckedOut struct {
	Umwelt       *umwelt.Umwelt
	Options      checkout_store.CheckinOptions
	AllowMissing bool
}

type ReadCheckedOutResults struct {
	Zettelen map[hinweis.Hinweis]stored_zettel.CheckedOut
}

func (op ReadCheckedOut) RunOneHinweis(
	s store_with_lock.Store,
	h hinweis.Hinweis,
) (zettel stored_zettel.CheckedOut, err error) {
	return op.RunOneString(s, h.String())
}

func (op ReadCheckedOut) RunOneString(
	s store_with_lock.Store,
	path string,
) (zettel stored_zettel.CheckedOut, err error) {
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
	results.Zettelen = make(map[hinweis.Hinweis]stored_zettel.CheckedOut)
	for _, p := range paths {
		var checked_out stored_zettel.CheckedOut

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
) (zettel stored_zettel.CheckedOut, err error) {
	if op.Options.AddMdExtension {
		p = p + ".md"
	}

	zettel.External, err = store.CheckoutStore().Read(p)

	if op.Options.IgnoreMissingHinweis && errors.IsNotExist(err) {
		err = nil
		//results.Zettelen[ez.Hinweis] = stored_zettel.External{}
		// continue
	} else if err != nil {
		err = errors.Error(err)
		return
	}

	if zettel.Internal, err = store.Zettels().Read(zettel.External.Hinweis); err != nil {
		err = errors.Wrapped(err, "%s", zettel.External.Path)
		return
	}

	return
}
