package query

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type VirtualStoreInitable struct {
	VirtualStore
	didInit  bool
	onceInit sync.Once
}

func (ve *VirtualStoreInitable) Initialize() (err error) {
	ve.onceInit.Do(func() {
		err = ve.VirtualStore.Initialize()
		ve.didInit = true
	})

	return
}

func (ve *VirtualStoreInitable) Flush() (err error) {
	if !ve.didInit {
		return
	}

	if err = ve.VirtualStore.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *VirtualStoreInitable) QueryCheckedOut(
	qg sku.ExternalQuery,
	f schnittstellen.FuncIter[sku.CheckedOutLike],
) (err error) {
	esqco, ok := es.VirtualStore.(sku.ExternalStoreQueryCheckedOut)

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

func (es *VirtualStoreInitable) CheckoutOne(
	options checkout_options.Options,
	sz *sku.Transacted,
) (cz sku.CheckedOutLike, err error) {
	escoo, ok := es.VirtualStore.(sku.ExternalStoreCheckoutOne)

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

type Virtual struct {
	sku.Queryable
	*Kennung
}

// func (ve *VirtualStoreInitable) Query(
// 	qg *Group,
// 	f schnittstellen.FuncIter[*sku.Transacted],
// ) (err error) {
// 	if err = ve.Initialize(); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if err = ve.VirtualStore.Query(qg, f); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

// func (ve *Virtual) ContainsSku(sk *sku.Transacted) bool {
// 	if !ve.Queryable.ContainsSku(sk) {
// 		return false
// 	}

// 	if !ve.Kennung.ContainsSku(sk) {
// 		return false
// 	}

// 	return true
// }
