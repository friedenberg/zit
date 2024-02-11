package matcher

import (
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
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
