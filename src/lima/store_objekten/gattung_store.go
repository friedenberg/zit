package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/foxtrot/sku"
)

type GattungStore interface {
}

type reindexer interface {
	reindexOne(sku.DataIdentity) (schnittstellen.Stored, error)
}