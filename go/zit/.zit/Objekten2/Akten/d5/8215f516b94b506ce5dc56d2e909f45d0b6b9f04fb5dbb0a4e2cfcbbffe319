package external_store

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

type Store struct {
	Supplies
	StoreLike

	didInit   bool
	onceInit  sync.Once
	initError error
}

func (ve *Store) Initialize() (err error) {
	ve.onceInit.Do(func() {
		ve.initError = ve.StoreLike.Initialize(ve.Supplies)
		ve.didInit = true
	})

	err = ve.initError

	return
}

func (ve *Store) Flush() (err error) {
	if !ve.didInit {
		return
	}

	if err = ve.StoreLike.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryCheckedOut(
	qg *query.Group,
	f interfaces.FuncIter[sku.SkuType],
) (err error) {
	es, ok := s.StoreLike.(QueryCheckedOut)

	if !ok {
		err = makeErrUnsupportedOperation(s, &s)
		return
	}

	if err = s.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = es.QueryCheckedOut(
		qg,
		f,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *Store) ReadAllExternalItems() (err error) {
	esado, ok := es.StoreLike.(sku.ExternalStoreReadAllExternalItems)

	if !ok {
		err = errors.Errorf("store does not support %T", &esado)
		return
	}

	if err = es.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = esado.ReadAllExternalItems(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadTransactedFromObjectId(
	o sku.CommitOptions,
	k1 interfaces.ObjectId,
	t *sku.Transacted,
) (e sku.ExternalLike, err error) {
	es, ok := s.StoreLike.(sku.ExternalStoreReadExternalLikeFromObjectId)

	if !ok {
		err = makeErrUnsupportedOperation(s, &s)
		return
	}

	if err = s.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if e, err = es.ReadExternalLikeFromObjectId(
		o,
		k1,
		t,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadExternalLikeFromObjectId(
	o sku.CommitOptions,
	k1 interfaces.ObjectId,
	t *sku.Transacted,
) (e sku.ExternalLike, err error) {
	es, ok := s.StoreLike.(sku.ExternalStoreReadExternalLikeFromObjectId)

	if !ok {
		err = makeErrUnsupportedOperation(s, &s)
		return
	}

	if err = s.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if e, err = es.ReadExternalLikeFromObjectId(
		o,
		k1,
		t,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) CheckoutOne(
	options checkout_options.Options,
	sz sku.TransactedGetter,
) (cz sku.SkuType, err error) {
	es, ok := s.StoreLike.(CheckoutOne)

	if !ok {
		err = makeErrUnsupportedOperation(s, &s)
		return
	}

	if err = s.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if cz, err = es.CheckoutOne(
		options,
		sz,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) DeleteCheckedOut(el *sku.CheckedOut) (err error) {
	es, ok := s.StoreLike.(DeleteCheckedOut)

	if !ok {
		err = makeErrUnsupportedOperation(s, &s)
		return
	}

	if err = s.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = es.DeleteCheckedOut(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// Takes a given sku.Transacted (internal) and updates it with the state of its
// checked out counterpart (external).
func (s *Store) UpdateTransacted(z *sku.Transacted) (err error) {
	es, ok := s.StoreLike.(UpdateTransacted)

	if !ok {
		err = makeErrUnsupportedOperation(s, &s)
		return
	}

	if err = s.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = es.UpdateTransacted(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) UpdateTransactedFromBlobs(z sku.ExternalLike) (err error) {
	es, ok := s.StoreLike.(UpdateTransactedFromBlobs)

	if !ok {
		err = makeErrUnsupportedOperation(s, &es)
		return
	}

	if err = s.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = es.UpdateTransactedFromBlobs(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *Store) GetObjectIdsForString(v string) (k []sku.ExternalObjectId, err error) {
	if es == nil {
		err = collections.MakeErrNotFoundString(v)
		return
	}

	if err = es.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if k, err = es.StoreLike.GetObjectIdsForString(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) Open(
	m checkout_mode.Mode,
	ph interfaces.FuncIter[string],
	zsc sku.SkuTypeSet,
) (err error) {
	es, ok := s.StoreLike.(Open)

	if !ok {
		err = makeErrUnsupportedOperation(s, &es)
		return
	}

	if err = s.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = es.Open(m, ph, zsc); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) SaveBlob(el sku.ExternalLike) (err error) {
	es, ok := s.StoreLike.(sku.BlobSaver)

	if !ok {
		err = makeErrUnsupportedOperation(s, &es)
		return
	}

	if err = s.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = es.SaveBlob(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) UpdateCheckoutFromCheckedOut(
	options checkout_options.OptionsWithoutMode,
	col sku.SkuType,
) (err error) {
	es, ok := s.StoreLike.(UpdateCheckoutFromCheckedOut)

	if !ok {
		err = makeErrUnsupportedOperation(s, &s)
		return
	}

	if err = s.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = es.UpdateCheckoutFromCheckedOut(options, col); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadCheckedOutFromTransacted(
	sk *sku.Transacted,
) (co *sku.CheckedOut, err error) {
	es, ok := s.StoreLike.(ReadCheckedOutFromTransacted)

	if !ok {
		err = makeErrUnsupportedOperation(s, &s)
		return
	}

	if err = s.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if co, err = es.ReadCheckedOutFromTransacted(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
