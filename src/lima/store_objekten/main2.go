package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/transacted"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

func (s *Store) ReadOne(
	k kennung.KennungPtr,
) (sk sku.SkuLikePtr, err error) {
	switch kt := k.(type) {
	case *kennung.Hinweis:
		return s.Zettel().ReadOne(kt)

	case *kennung.Typ:
		return s.typStore.ReadOne(kt)

	case *kennung.Etikett:
		return s.Etikett().ReadOne(kt)

	case *kennung.Kasten:
		return s.Kasten().ReadOne(kt)

	case *kennung.Konfig:
		return s.Konfig().ReadOne(kt)

	default:
		err = errors.Errorf("unsupported kennung %T -> %q", kt, kt)
		return
	}
}

func (s *Store) ReadAllSchwanzen(
	gs gattungen.Set,
	f schnittstellen.FuncIter[sku.SkuLikePtr],
) (err error) {
	chErr := make(chan error, gs.Len())

	for g, s1 := range s.readers {
		if !gs.ContainsKey(g.GetGattungString()) {
			continue
		}

		go func(s1 objekte.FuncReaderTransactedLikePtr) {
			var subErr error

			defer func() {
				chErr <- subErr
			}()

			subErr = s1(f)
		}(s1)
	}

	for i := 0; i < gs.Len(); i++ {
		err = errors.MakeMulti(err, <-chErr)
	}

	return
}

func (s *Store) ReadAll(
	gs gattungen.Set,
	f schnittstellen.FuncIter[sku.SkuLikePtr],
) (err error) {
	chErr := make(chan error, gs.Len())

	for g, s1 := range s.transactedReaders {
		if !gs.ContainsKey(g.GetGattungString()) {
			continue
		}

		go func(s1 objekte.FuncReaderTransactedLikePtr) {
			var subErr error

			defer func() {
				chErr <- subErr
			}()

			subErr = s1(f)
		}(s1)
	}

	for i := 0; i < gs.Len(); i++ {
		err = errors.MakeMulti(err, <-chErr)
	}

	return
}

func (s *Store) CheckoutOne(
	options CheckoutOptions,
	sk sku.SkuLikePtr,
) (co objekte.CheckedOutLikePtr, err error) {
	switch skt := sk.(type) {
	case *transacted.Zettel:
		return s.Zettel().CheckoutOne(options, skt)

	case *transacted.Typ:
		return s.typStore.CheckoutOne(options, skt)

	case *transacted.Etikett:
		return s.Etikett().CheckoutOne(options, skt)

	case *transacted.Kasten:
		return s.Kasten().CheckoutOne(options, skt)

	case *transacted.Konfig:
		return s.Konfig().CheckoutOne(options, skt)

	default:
		err = errors.Errorf("unsupported kennung %T -> %q", skt, skt)
		return
	}
}
