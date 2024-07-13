package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
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
		interfaces.Stringer
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
			f interfaces.FuncIter[*Transacted],
			k kennung.Kasten,
			updateTransacted func(
				kasten kennung.Kasten,
				z *Transacted,
			) (err error),
		) interfaces.FuncIter[*Transacted]
		MakeEmitSkuSigilSchwanzen(
			f interfaces.FuncIter[*Transacted],
			k kennung.Kasten,
			updateTransacted func(
				kasten kennung.Kasten,
				z *Transacted,
			) (err error),
		) interfaces.FuncIter[*Transacted]
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
