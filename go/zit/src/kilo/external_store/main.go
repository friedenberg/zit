package external_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
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
		) (cz sku.SkuType, err error)
	}

	DeleteCheckedOut interface {
		DeleteCheckedOut(el *sku.CheckedOut) (err error)
	}

	UpdateTransacted = sku.ExternalStoreUpdateTransacted

	UpdateTransactedFromBlobs interface {
		UpdateTransactedFromBlobs(sku.ExternalLike) (err error)
	}

	Open interface {
		Open(
			m checkout_mode.Mode,
			ph interfaces.FuncIter[string],
			zsc sku.SkuTypeSet,
		) (err error)
	}

	UpdateCheckoutFromCheckedOut interface {
		UpdateCheckoutFromCheckedOut(
			options checkout_options.OptionsWithoutMode,
			co sku.SkuType,
		) (err error)
	}

	ReadCheckedOutFromTransacted interface {
		ReadCheckedOutFromTransacted(
			sk *sku.Transacted,
		) (co *sku.CheckedOut, err error)
	}

	QueryCheckedOut = query.QueryCheckedOut

	Supplies struct {
		StoreFuncs
		DirCache string
		dir_layout.DirLayout
		ids.RepoId
		ids.TypeSet
		ids.Clock
		BlobStore *blob_store.VersionedStores
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
