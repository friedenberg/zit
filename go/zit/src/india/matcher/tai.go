package matcher

import (
	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/echo/kennung"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
)

type Tai kennung.Tai

func (t Tai) ContainsMatchableExactly(m *sku.Transacted) bool {
	return false
}

func (t Tai) ContainsMatchable(m *sku.Transacted) bool {
	errors.TodoP1("add GetTai to matchable")
	return false
}
