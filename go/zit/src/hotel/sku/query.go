package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type (
	Queryable interface {
		// AppendMatchToQueryPath(*Transacted, *QueryPath) error
		ContainsSku(TransactedGetter) bool
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

	FuncQuery = func(
		QueryGroup,
		interfaces.FuncIter[*Transacted],
	) (err error)

	FuncPrimitiveQuery = func(
		PrimitiveQueryGroup,
		interfaces.FuncIter[*Transacted],
	) (err error)

	QueryGroup interface {
		PrimitiveQueryGroup
		Query
		SetIncludeHistory()
		GetTags() ids.TagSet
		GetTypes() ids.TypeSet
	}

	ExternalQueryOptions struct {
		ids.RepoId
		ExcludeUntracked  bool
		ExcludeRecognized bool
	}

	ExternalQuery struct {
		ids.RepoId
		QueryGroup
		ExternalQueryOptions
	}
)
