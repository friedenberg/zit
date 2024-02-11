package matcher

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type Tai kennung.Tai

func (t Tai) ContainsMatchableExactly(m *sku.Transacted) bool {
	return false
}

func (t Tai) ContainsMatchable(m *sku.Transacted) bool {
	errors.TodoP1("add GetTai to matchable")
	return false
}
