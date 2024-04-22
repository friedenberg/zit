package objekte

import (
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/india/query"
)

type (
	Konfig interface {
		schnittstellen.Konfig
		query.ImplicitEtikettenGetter
		IsInlineTyp(kennung.Typ) bool
		GetApproximatedTyp(kennung.Kennung) ApproximatedTyp
	}
)
