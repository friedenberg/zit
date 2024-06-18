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

	QueryWithSigilAndKennung interface {
		Query
		GetSigil() kennung.Sigil
		ContainsKennung(*kennung.Kennung2) bool
	}

	QueryGroup interface {
		Query
		HasHidden() bool
		Get(gattung.Gattung) (QueryWithSigilAndKennung, bool)
	}
)
