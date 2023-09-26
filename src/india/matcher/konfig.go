package matcher

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
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
