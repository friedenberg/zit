package matcher

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type Tai kennung.Tai

func (t Tai) ContainsMatchableExactly(m *sku.Transacted) bool {
	return false
}

func (t Tai) ContainsMatchable(m *sku.Transacted) bool {
	errors.TodoP1("add GetTai to matchable")
	return false
}
