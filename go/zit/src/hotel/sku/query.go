package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type (
	Queryable interface {
		// AppendMatchToQueryPath(*Transacted, *QueryPath) error
		ContainsSku(*Transacted) bool
	}

	Query interface {
		Queryable
		interfaces.Stringer
		// Each(schnittstellen.FuncIter[Query]) error
	}

	SigilGetter interface {
		GetSigil() ids.Sigil
	}

	QueryWithSigilAndObjectId interface {
		Query
		SigilGetter
		ContainsObjectId(*ids.ObjectId) bool
	}

	// Used by store_verzeichnisse.binary*
	PrimitiveQueryGroup interface {
		Get(genres.Genre) (QueryWithSigilAndObjectId, bool)
		SigilGetter
		HasHidden() bool
	}

	QueryGroup interface {
		PrimitiveQueryGroup
		Query
		SetIncludeHistory()
		MakeEmitSkuMaybeExternal(
			f interfaces.FuncIter[*Transacted],
			k ids.RepoId,
			updateTransacted func(
				kasten ids.RepoId,
				z *Transacted,
			) (err error),
		) interfaces.FuncIter[*Transacted]
		MakeEmitSkuSigilLatest(
			f interfaces.FuncIter[*Transacted],
			k ids.RepoId,
			updateTransacted func(
				kasten ids.RepoId,
				z *Transacted,
			) (err error),
		) interfaces.FuncIter[*Transacted]
		GetTags() ids.TagSet
		GetTypes() ids.TypeSet
	}

	ExternalQueryOptions struct {
		ids.RepoId
		ExcludeUntracked  bool
		IncludeRecognized bool
	}

	ExternalQuery struct {
		ids.RepoId
		QueryGroup
		ExternalQueryOptions
	}
)
