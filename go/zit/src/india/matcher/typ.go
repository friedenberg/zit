package matcher

import (
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type Typ kennung.Typ

func (t Typ) ContainsMatchableExactly(m *sku.Transacted) bool {
	g := gattung.Make(m.GetGattung())

	switch g {
	case gattung.Zettel, gattung.Typ:
		// noop
	default:
		return false
	}

	t1 := m.GetTyp()

	if kennung.Typ(t).Equals(t1) {
		return true
	}

	t2, ok := m.GetKennung().(kennung.Typ)

	if !ok {
		return false
	}

	if !kennung.Typ(t).Equals(t2) {
		return false
	}

	return true
}
