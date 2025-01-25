package external_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
)

type (
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

	MergeCheckedOut interface {
		MergeCheckedOut(
			co *sku.CheckedOut,
			parentNegotiator sku.ParentNegotiator,
			allowMergeConflicts bool,
		) (commitOptions sku.CommitOptions, err error)
	}

	QueryCheckedOut = query.QueryCheckedOut

	Supplies struct {
		sku.ObjectStore
		DirCache string
		env_repo.Env
		ids.RepoId
		ids.TypeSet
		ids.Clock
		BlobStore *typed_blob_store.Store
	}

	StoreLike interface {
		Initialize(Supplies) error
		QueryCheckedOut
		errors.Flusher
		sku.ExternalStoreForQuery
	}

	StoreGetter interface {
		GetExternalStore(ids.RepoId) (*Store, bool)
	}
)
