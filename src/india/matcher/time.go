package matcher

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type Time kennung.Time

func (t Time) ContainsMatchable(m *sku.Transacted) bool {
	return false
}
