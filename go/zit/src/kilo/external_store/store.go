package external_store

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
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

type (
	FuncRealize     = func(*Transacted, *Transacted, ObjekteOptions) error
	FuncCommit      = func(*Transacted, ObjekteOptions) error
	FuncReadSha     = func(*sha.Sha) (*Transacted, error)
	FuncReadOneInto = func(
		k1 interfaces.StringerGenreGetter,
		out *Transacted,
	) (err error)

	StoreFuncs struct {
		FuncRealize
		FuncCommit
		FuncReadSha
		FuncReadOneInto
		sku.FuncQuery
	}

	QueryOptions struct {
		ExcludeUntracked  bool
		IncludeRecognized bool
	}

	CheckoutOne interface {
		CheckoutOne(
			options checkout_options.Options,
			sz *Transacted,
		) (cz CheckedOutLike, err error)
	}

	DeleteCheckout interface {
		DeleteCheckout(col CheckedOutLike) (err error)
	}

	UpdateTransacted interface {
		UpdateTransacted(z *Transacted) (err error)
	}

	Open interface {
		Open(
			m checkout_mode.Mode,
			ph interfaces.FuncIter[string],
			zsc CheckedOutLikeSet,
		) (err error)
	}

	QueryCheckedOut interface {
		QueryCheckedOut(
			qg *query.Group,
			f interfaces.FuncIter[CheckedOutLike],
		) (err error)
	}

	QueryUnsure interface {
		QueryUnsure(
			qg *query.Group,
			f interfaces.FuncIter[CheckedOutLike],
		) (err error)
	}

	Info struct {
		StoreFuncs
		DirCache string
		standort.Standort
	}

	StoreLike interface {
		Initialize(Info) error
		QueryUnsure
		QueryCheckedOut
		interfaces.Flusher
		sku.ExternalStoreForQuery
	}

	StoreGetter interface {
		GetExternalStore(kennung.RepoId) (*Store, bool)
	}
)

// Add typ set
type Store struct {
	kennung.TypSet
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
	f interfaces.FuncIter[CheckedOutLike],
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

func (es *Store) QueryUnsure(
	qg *query.Group,
	f interfaces.FuncIter[CheckedOutLike],
) (err error) {
	esqu, ok := es.StoreLike.(QueryUnsure)

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

func (es *Store) CheckoutOne(
	options checkout_options.Options,
	sz *Transacted,
) (cz CheckedOutLike, err error) {
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

func (es *Store) DeleteCheckout(col CheckedOutLike) (err error) {
	esdc, ok := es.StoreLike.(DeleteCheckout)

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

func (es *Store) UpdateTransacted(z *Transacted) (err error) {
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

func (es *Store) GetExternalKennung() (ks interfaces.SetLike[*kennung.Id], err error) {
	if es == nil {
		ks = collections_value.MakeValueSet[*kennung.Id](nil)
		return
	}

	if err = es.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if ks, err = es.StoreLike.GetExternalKennung(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *Store) GetKennungForString(v string) (k *kennung.Id, err error) {
	if es == nil {
		err = collections.MakeErrNotFoundString(v)
		return
	}

	if err = es.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if k, err = es.StoreLike.GetKennungForString(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es *Store) Open(
	m checkout_mode.Mode,
	ph interfaces.FuncIter[string],
	zsc CheckedOutLikeSet,
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

type ErrExternalStoreUnsupportedTyp kennung.Typ

func (e ErrExternalStoreUnsupportedTyp) Is(target error) bool {
	_, ok := target.(ErrExternalStoreUnsupportedTyp)
	return ok
}

func (e ErrExternalStoreUnsupportedTyp) Error() string {
	return fmt.Sprintf("unsupported typ: %q", kennung.Typ(e))
}
