package external_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

type (
	FuncRealize     = func(sku.ExternalLike, *sku.Transacted, sku.CommitOptions) error
	FuncCommit      = func(sku.ExternalLike, sku.CommitOptions) error
	FuncReadOneInto = func(string, *sku.Transacted) error

	StoreFuncs struct {
		FuncRealize
		FuncCommit
		FuncReadOneInto
		sku.FuncPrimitiveQuery
	}

	QueryOptions struct {
		ExcludeUntracked  bool
		IncludeRecognized bool
	}

	CheckoutOne interface {
		CheckoutOne(
			options checkout_options.Options,
			sz sku.TransactedGetter,
		) (cz sku.CheckedOutLike, err error)
	}

	DeleteExternal interface {
		DeleteExternalLike(el sku.ExternalLike) (err error)
	}

	UpdateTransacted = sku.ExternalStoreUpdateTransacted

	Open interface {
		Open(
			m checkout_mode.Mode,
			ph interfaces.FuncIter[string],
			zsc sku.CheckedOutLikeSet,
		) (err error)
	}

	UpdateCheckoutFromCheckedOut interface {
		UpdateCheckoutFromCheckedOut(
			options checkout_options.OptionsWithoutMode,
			co sku.CheckedOutLike,
		) (err error)
	}

	QueryCheckedOut = query.QueryCheckedOut

	Supplies struct {
		StoreFuncs
		DirCache string
		fs_home.Home
		ids.RepoId
		ids.TypeSet
		ids.Clock
	}

	StoreLike interface {
		Initialize(Supplies) error
		QueryCheckedOut
		interfaces.Flusher
		sku.ExternalStoreForQuery
	}

	StoreGetter interface {
		GetExternalStore(ids.RepoId) (*Store, bool)
	}
)
