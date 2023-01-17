package store_objekten

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/golf/sku"
)

type GattungStore interface {
}

type reindexer interface {
	reindexOne(sku.DataIdentity) (gattung.Stored, error)
}
