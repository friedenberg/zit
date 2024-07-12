package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

type (
	Queryable interface {
		// AppendMatchToQueryPath(*Transacted, *QueryPath) error
		ContainsSku(*Transacted) bool
	}

	Query interface {
		Queryable
		schnittstellen.Stringer
		// Each(schnittstellen.FuncIter[Query]) error
	}

	SigilGetter interface {
		GetSigil() kennung.Sigil
	}

	QueryWithSigilAndKennung interface {
		Query
		SigilGetter
		ContainsKennung(*kennung.Kennung2) bool
	}

	// Used by store_verzeichnisse.binary*
	PrimitiveQueryGroup interface {
		Get(gattung.Gattung) (QueryWithSigilAndKennung, bool)
		SigilGetter
		HasHidden() bool
	}

	QueryGroup interface {
		PrimitiveQueryGroup
		Query
		SetIncludeHistory()
		MakeEmitSkuMaybeExternal(
			f schnittstellen.FuncIter[*Transacted],
			k kennung.Kasten,
			updateTransacted func(
				kasten kennung.Kasten,
				z *Transacted,
			) (err error),
		) schnittstellen.FuncIter[*Transacted]
		MakeEmitSkuSigilSchwanzen(
			f schnittstellen.FuncIter[*Transacted],
			k kennung.Kasten,
			updateTransacted func(
				kasten kennung.Kasten,
				z *Transacted,
			) (err error),
		) schnittstellen.FuncIter[*Transacted]
		GetEtiketten() kennung.EtikettSet
		GetTypen() kennung.TypSet
	}

	ExternalQueryOptions struct {
		kennung.Kasten
		ExcludeUntracked  bool
		IncludeRecognized bool
	}

	ExternalQuery struct {
		kennung.Kasten
		QueryGroup
		ExternalQueryOptions
	}
)
