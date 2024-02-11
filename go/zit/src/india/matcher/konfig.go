package matcher

import (
	"code.linenisgreat.com/zit-go/src/charlie/gattung"
	"code.linenisgreat.com/zit-go/src/echo/kennung"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
)

type Konfig kennung.Konfig

func (i Konfig) ContainsMatchable(m *sku.Transacted) bool {
	if !kennung.Konfig(
		i,
	).GetGattung().EqualsGattung(
		gattung.Make(m.GetGattung()),
	) {
		return false
	}

	return true
}
