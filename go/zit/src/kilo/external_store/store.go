package external_store

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

type Store struct {
	Info
	StoreLike

	didInit  bool
	onceInit sync.Once
}

func (ve *Store) Initialize() (err error) {
	ve.onceInit.Do(func() {
		err = ve.StoreLike.Initialize(ve.Info)
		ve.didInit = true
	})

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

func (es *Store) QueryCheckedOut(
	qg *query.Group,
	f interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	esqco, ok := es.StoreLike.(QueryCheckedOut)

	if !ok {
		err = errors.Errorf("store does not support %T", esqco)
		return
	}

	if err = es.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = esqco.QueryCheckedOut(
		qg,
		f,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *Store) ReadTransactedFromObjectId(
	o sku.CommitOptions,
	k1 interfaces.ObjectId,
	t *sku.Transacted,
) (e sku.ExternalLike, err error) {
	esrtfoi, ok := es.StoreLike.(ReadTransactedFromObjectId)

	if !ok {
		err = errors.Errorf("store does not support %T", esrtfoi)
		return
	}

	if err = es.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if e, err = esrtfoi.ReadTransactedFromObjectId(
		o,
		k1,
		t,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *Store) CheckoutOne(
	options checkout_options.Options,
	sz *sku.Transacted,
) (cz sku.CheckedOutLike, err error) {
	escoo, ok := es.StoreLike.(CheckoutOne)

	if !ok {
		err = errors.Errorf("store does not support %T", escoo)
		return
	}

	if err = es.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if cz, err = escoo.CheckoutOne(
		options,
		sz,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *Store) DeleteCheckout(col sku.CheckedOutLike) (err error) {
	esdc, ok := es.StoreLike.(DeleteExternal)

	if !ok {
		err = errors.Errorf("store does not support DeleteCheckout")
		return
	}

	if err = es.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = esdc.DeleteCheckout(col); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *Store) UpdateTransacted(z *sku.Transacted) (err error) {
	esut, ok := es.StoreLike.(UpdateTransacted)

	if !ok {
		err = errors.Errorf("store does not support UpdateTransacted")
		return
	}

	if err = es.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = esut.UpdateTransacted(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *Store) GetExternalObjectIds() (ks interfaces.SetLike[*ids.ObjectId], err error) {
	if es == nil {
		ks = collections_value.MakeValueSet[*ids.ObjectId](nil)
		return
	}

	if err = es.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if ks, err = es.StoreLike.GetExternalObjectIds(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *Store) GetObjectsIdForString(v string) (k []*ids.ObjectId, err error) {
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

func (es *Store) Open(
	m checkout_mode.Mode,
	ph interfaces.FuncIter[string],
	zsc sku.CheckedOutLikeSet,
) (err error) {
	eso, ok := es.StoreLike.(Open)

	if !ok {
		err = errors.Errorf("store does not support UpdateTransacted")
		return
	}

	if err = es.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = eso.Open(m, ph, zsc); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) GetExternalStoreOrganizeFormat(
	f *sku_fmt.Organize,
) sku_fmt.ExternalLike {
	esof, ok := s.StoreLike.(OrganizeFormatGetter)

	if !ok {
		return f
	}

	return esof.GetExternalStoreOrganizeFormat(f)
}
