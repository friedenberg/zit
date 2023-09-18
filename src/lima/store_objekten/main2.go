package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/transacted"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

func (s *Store) onNewOrUpdated(
	t *sku.Transacted2,
) (err error) {
	s.StoreUtil.CommitUpdatedTransacted(t)

	g := gattung.Must(t.Kennung.GetGattung())

	switch g {
	case gattung.Typ:
		if err = s.StoreUtil.GetKonfigPtr().AddTyp(t); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = gattung.MakeErrUnsupportedGattung(g)
		return
	}

	return
}

func (s *Store) onNew(
	t *sku.Transacted2,
) (err error) {
	if err = s.onNewOrUpdated(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return s.LogWriter.New(t)
}

func (s *Store) onUpdated(
	t *sku.Transacted2,
) (err error) {
	if err = s.onNewOrUpdated(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return s.LogWriter.Updated(t)
}

func (s *Store) onUnchanged(
	t *sku.Transacted2,
) (err error) {
	return s.LogWriter.Unchanged(t)
}

func (s *Store) ReadOne(
	k *kennung.Kennung2,
) (sk *sku.Transacted2, err error) {
	switch kt := k.KennungPtr.(type) {
	// case *kennung.Hinweis:
	// 	return s.Zettel().ReadOne(kt)

	case *kennung.Typ:
		var sk1 *transacted.Typ

		if sk1, err = s.typStore.ReadOne(kt); err != nil {
			err = errors.Wrap(err)
			return
		}

		sk = &sku.Transacted2{}

		if err = sk.SetFromSkuLike(sk1); err != nil {
			err = errors.Wrap(err)
			return
		}

	// case *kennung.Etikett:
	// 	return s.Etikett().ReadOne(kt)

	// case *kennung.Kasten:
	// 	return s.Kasten().ReadOne(kt)

	// case *kennung.Konfig:
	// 	return s.Konfig().ReadOne(kt)

	default:
		err = errors.Errorf("unsupported kennung %T -> %q", kt, kt)
		return
	}

	return
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
