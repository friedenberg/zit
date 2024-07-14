package sku

import (
	"fmt"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type (
	FuncRealize     = func(*Transacted, *Transacted, ObjectOptions) error
	FuncCommit      = func(*Transacted, ObjectOptions) error
	FuncReadSha     = func(*sha.Sha) (*Transacted, error)
	FuncReadOneInto = func(
		k1 interfaces.ObjectId,
		out *Transacted,
	) (err error)

	FuncQuery = func(
		QueryGroup,
		interfaces.FuncIter[*Transacted],
	) (err error)

	StoreFuncs struct {
		FuncRealize
		FuncCommit
		FuncReadSha
		FuncReadOneInto
		FuncQuery
	}

	ExternalStoreCheckoutOne interface {
		CheckoutOne(
			options checkout_options.Options,
			sz *Transacted,
		) (cz CheckedOutLike, err error)
	}

	ExternalStoreDeleteCheckout interface {
		DeleteCheckout(col CheckedOutLike) (err error)
	}

	ExternalStoreUpdateTransacted interface {
		UpdateTransacted(z *Transacted) (err error)
	}

	ExternalStoreOpen interface {
		Open(
			m checkout_mode.Mode,
			ph interfaces.FuncIter[string],
			zsc CheckedOutLikeSet,
		) (err error)
	}

	ExternalStoreQueryCheckedOut interface {
		QueryCheckedOut(
			qg ExternalQuery,
			f interfaces.FuncIter[CheckedOutLike],
		) (err error)
	}

	ExternalStoreQueryUnsure interface {
		QueryUnsure(
			qg ExternalQuery,
			f interfaces.FuncIter[CheckedOutLike],
		) (err error)
	}

	ExternalStoreInfo struct {
		StoreFuncs
		DirCache string
		fs_home.Home
	}

	ExternalStoreForQuery interface {
		GetExternalKennung() (interfaces.SetLike[*ids.ObjectId], error)
		GetKennungForString(string) (*ids.ObjectId, error)
	}

	ExternalStoreLike interface {
		Initialize(ExternalStoreInfo) error
		ExternalStoreQueryUnsure
		ExternalStoreQueryCheckedOut
		interfaces.Flusher
		ExternalStoreForQuery
	}

	ExternalStoreForQueryGetter interface {
		GetExternalStoreForQuery(ids.RepoId) (ExternalStoreForQuery, bool)
	}
)

// Add typ set
type ExternalStore struct {
	ids.TypeSet
	ExternalStoreInfo
	ExternalStoreLike
	didInit  bool
	onceInit sync.Once
}

func (ve *ExternalStore) Initialize() (err error) {
	ve.onceInit.Do(func() {
		err = ve.ExternalStoreLike.Initialize(ve.ExternalStoreInfo)
		ve.didInit = true
	})

	return
}

func (ve *ExternalStore) Flush() (err error) {
	if !ve.didInit {
		return
	}

	if err = ve.ExternalStoreLike.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *ExternalStore) QueryCheckedOut(
	qg ExternalQuery,
	f interfaces.FuncIter[CheckedOutLike],
) (err error) {
	esqco, ok := es.ExternalStoreLike.(ExternalStoreQueryCheckedOut)

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

func (es *ExternalStore) QueryUnsure(
	qg ExternalQuery,
	f interfaces.FuncIter[CheckedOutLike],
) (err error) {
	esqu, ok := es.ExternalStoreLike.(ExternalStoreQueryUnsure)

	if !ok {
		err = errors.Errorf("store does not support %T", esqu)
		return
	}

	if err = es.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = esqu.QueryUnsure(
		qg,
		f,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *ExternalStore) CheckoutOne(
	options checkout_options.Options,
	sz *Transacted,
) (cz CheckedOutLike, err error) {
	escoo, ok := es.ExternalStoreLike.(ExternalStoreCheckoutOne)

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

func (es *ExternalStore) DeleteCheckout(col CheckedOutLike) (err error) {
	esdc, ok := es.ExternalStoreLike.(ExternalStoreDeleteCheckout)

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

func (es *ExternalStore) UpdateTransacted(z *Transacted) (err error) {
	esut, ok := es.ExternalStoreLike.(ExternalStoreUpdateTransacted)

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

func (es *ExternalStore) GetExternalKennung() (ks interfaces.SetLike[*ids.ObjectId], err error) {
	if es == nil {
		ks = collections_value.MakeValueSet[*ids.ObjectId](nil)
		return
	}

	if err = es.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if ks, err = es.ExternalStoreLike.GetExternalKennung(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *ExternalStore) GetKennungForString(v string) (k *ids.ObjectId, err error) {
	if es == nil {
		err = collections.MakeErrNotFoundString(v)
		return
	}

	if err = es.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if k, err = es.ExternalStoreLike.GetKennungForString(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *ExternalStore) Open(
	m checkout_mode.Mode,
	ph interfaces.FuncIter[string],
	zsc CheckedOutLikeSet,
) (err error) {
	eso, ok := es.ExternalStoreLike.(ExternalStoreOpen)

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

type ErrExternalStoreUnsupportedTyp ids.Type

func (e ErrExternalStoreUnsupportedTyp) Is(target error) bool {
	_, ok := target.(ErrExternalStoreUnsupportedTyp)
	return ok
}

func (e ErrExternalStoreUnsupportedTyp) Error() string {
	return fmt.Sprintf("unsupported typ: %q", ids.Type(e))
}
