package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
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

	QueryWithSigilAndKennung interface {
		Query
		SigilGetter
		ContainsKennung(*ids.ObjectId) bool
	}

	// Used by store_verzeichnisse.binary*
	PrimitiveQueryGroup interface {
		Get(gattung.Genre) (QueryWithSigilAndKennung, bool)
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
		MakeEmitSkuSigilSchwanzen(
			f interfaces.FuncIter[*Transacted],
			k ids.RepoId,
			updateTransacted func(
				kasten ids.RepoId,
				z *Transacted,
			) (err error),
		) interfaces.FuncIter[*Transacted]
		GetEtiketten() ids.TagSet
		GetTypen() ids.TypeSet
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
