package sku

import (
	"fmt"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

type (
	ExternalQuery struct {
		Queryable
		ExcludeUntracked bool
	}

	ExternalQueryWithKasten struct {
		ExternalQuery
		kennung.Kasten
	}

	FuncRealize     = func(*Transacted, *Transacted, ObjekteOptions) error
	FuncCommit      = func(*Transacted, ObjekteOptions) error
	FuncReadSha     = func(*sha.Sha) (*Transacted, error)
	FuncReadOneInto = func(
		k1 schnittstellen.StringerGattungGetter,
		out *Transacted,
	) (err error)

	FuncQuery = func(
		QueryGroup,
		schnittstellen.FuncIter[*Transacted],
	) (err error)

	StoreFuncs struct {
		FuncRealize
		FuncCommit
		FuncReadSha
		FuncReadOneInto
		FuncQuery
	}

	ExternalStoreQueryCheckedOut interface {
		QueryCheckedOut(
			qg ExternalQuery,
			f schnittstellen.FuncIter[CheckedOutLike],
		) (err error)
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

	ExternalStoreInitInfo struct {
		StoreFuncs
	}

	ExternalStoreLike interface {
		Initialize(ExternalStoreInitInfo) error
		ExternalStoreQueryCheckedOut
		// ExternalStoreCheckoutOne
		schnittstellen.Flusher
	}
)

// Add typ set
type ExternalStore struct {
	kennung.TypSet
	ExternalStoreInitInfo
	ExternalStoreLike
	didInit  bool
	onceInit sync.Once
}

func (ve *ExternalStore) Initialize() (err error) {
	ve.onceInit.Do(func() {
		err = ve.ExternalStoreLike.Initialize(ve.ExternalStoreInitInfo)
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
	f schnittstellen.FuncIter[CheckedOutLike],
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
		err = errors.Errorf("store does not support %T", esdc)
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

type ErrExternalStoreUnsupportedTyp kennung.Typ

func (e ErrExternalStoreUnsupportedTyp) Is(target error) bool {
	_, ok := target.(ErrExternalStoreUnsupportedTyp)
	return ok
}

func (e ErrExternalStoreUnsupportedTyp) Error() string {
	return fmt.Sprintf("unsupported typ: %q", kennung.Typ(e))
}
