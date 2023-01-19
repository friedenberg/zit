package store_objekten

import (
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/schnittstellen"
)

type GattungStore interface {
}

type reindexer interface {
	reindexOne(sku.DataIdentity) (schnittstellen.Stored, error)
}
