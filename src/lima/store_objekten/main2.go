package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/iter2"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/mike/store_util"
)

func (s *Store) onNewOrUpdated(
	t *sku.Transacted2,
) (err error) {
	s.StoreUtil.CommitUpdatedTransacted(t)

	g := gattung.Must(t.Kennung.GetGattung())

	switch g {
	case gattung.Typ:
		if err = s.StoreUtil.GetKonfigPtr().AddTyp2(t); err != nil {
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
	k1 schnittstellen.StringerGattungGetter,
) (sk *sku.Transacted2, err error) {
	var k kennung.Kennung2

	if err = k.SetWithGattung(k1.String(), k1.GetGattung()); err != nil {
		err = errors.Wrap(err)
		return
	}

	switch kt := k.KennungPtr.(type) {
	// case *kennung.Hinweis:
	// 	return s.Zettel().ReadOne(kt)

	case *kennung.Typ:
		var sk1 sku.SkuLikePtr

		if sk1, err = s.typStore.ReadOne(kt); err != nil {
			err = errors.Wrap(err)
			return
		}

		sk = &sku.Transacted2{}

		if err = sk.SetFromSkuLike(sk1); err != nil {
			err = errors.Wrap(err)
			return
		}

	case *kennung.Etikett:
		var sk1 sku.SkuLikePtr

		if sk1, err = s.etikettStore.ReadOne(kt); err != nil {
			err = errors.Wrap(err)
			return
		}

		sk = &sku.Transacted2{}

		if err = sk.SetFromSkuLike(sk1); err != nil {
			err = errors.Wrap(err)
			return
		}

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
	if err = iter2.Parallel(
		gs,
		func(g gattung.Gattung) (err error) {
			r, ok := s.readers[g]

			if !ok {
				return
			}

			return r(f)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadAll(
	gs gattungen.Set,
	f schnittstellen.FuncIter[sku.SkuLikePtr],
) (err error) {
	if err = iter2.Parallel(
		gs,
		func(g gattung.Gattung) (err error) {
			r, ok := s.transactedReaders[g]

			if !ok {
				return
			}

			return r(f)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) CheckoutOne(
	options store_util.CheckoutOptions,
	sk sku.SkuLikePtr,
) (co objekte.CheckedOutLikePtr, err error) {
	g := gattung.Must(sk)
	switch g {
	case gattung.Zettel:
		return s.Zettel().CheckoutOne(options, sk)

	case gattung.Typ:
		return s.typStore.CheckoutOne(options, sk)

	case gattung.Etikett:
		return s.Etikett().CheckoutOne(options, sk)

	case gattung.Kasten:
		return s.Kasten().CheckoutOne(options, sk)

	case gattung.Konfig:
		return s.Konfig().CheckoutOne(options, sk)

	default:
		err = gattung.MakeErrUnsupportedGattung(g)
		return
	}
}
