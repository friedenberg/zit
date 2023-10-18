package matcher

import (
	"github.com/friedenberg/zit/src/delta/thyme"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type Time thyme.Time

func (t Time) ContainsMatchable(m *sku.Transacted) bool {
	return false
}
