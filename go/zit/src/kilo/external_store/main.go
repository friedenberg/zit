package external_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
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
			sz *sku.Transacted,
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

	OrganizeFormatGetter interface {
		GetExternalStoreOrganizeFormat(
			*sku_fmt.Box,
		) sku_fmt.ExternalLike
	}

	QueryCheckedOut = query.QueryCheckedOut

	Info struct {
		StoreFuncs
		DirCache string
		fs_home.Home
		ids.RepoId
		ids.TypeSet
	}

	StoreLike interface {
		Initialize(Info) error
		QueryCheckedOut
		interfaces.Flusher
		sku.ExternalStoreForQuery
		sku.ExternalLikePoolGetter
		sku.ExternalLikeResetter3Getter
	}

	StoreGetter interface {
		GetExternalStore(ids.RepoId) (*Store, bool)
	}
)
