package sku

import (
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

	ExternalStore interface {
		schnittstellen.Flusher
		ExternalStoreQueryCheckedOut
		ExternalStoreCheckoutOne
	}
)
