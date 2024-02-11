package matcher

import (
	"code.linenisgreat.com/zit-go/src/delta/thyme"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
)

type Time thyme.Time

func (t Time) ContainsMatchable(m *sku.Transacted) bool {
	return false
}
