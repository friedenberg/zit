package user_ops

import (
	"os"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/foxtrot/hinweisen"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/hotel/zettels"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type ReadCheckedOut struct {
	Umwelt       *umwelt.Umwelt
	Options      zettels.CheckinOptions
	AllowMissing bool
}

type ReadCheckedOutResults struct {
	Zettelen map[hinweis.Hinweis]stored_zettel.CheckedOut
}

func (op ReadCheckedOut) RunOne(path string) (zettel stored_zettel.CheckedOut, err error) {
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

func (op ReadCheckedOut) Run(paths ...string) (results ReadCheckedOutResults, err error) {
	results.Zettelen = make(map[hinweis.Hinweis]stored_zettel.CheckedOut)

	var store store_with_lock.Store

	if store, err = store_with_lock.New(op.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	for _, p := range paths {
		var checked_out stored_zettel.CheckedOut

		if checked_out, err = op.runOne(store, p); err != nil {
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

func (op ReadCheckedOut) runOne(store store_with_lock.Store, p string) (zettel stored_zettel.CheckedOut, err error) {
	if op.Options.AddMdExtension {
		p = p + ".md"
	}

	zettel.External, err = store.CheckoutStore().Read(p)

	if op.Options.IgnoreMissingHinweis && errors.Is(err, os.ErrNotExist) {
		err = nil
		//results.Zettelen[ez.Hinweis] = stored_zettel.External{}
		// continue
	} else if err != nil {
		err = errors.Error(err)
		return
	}

	logz.Print(zettel.External.Hinweis)

	if zettel.Internal, err = store.Zettels().Read(zettel.External.Hinweis); err != nil {
		err = errors.Wrapped(err, "%s", zettel.External.Path)
		return
	}

	return
}
