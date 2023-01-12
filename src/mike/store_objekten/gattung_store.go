package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/golf/sku"
)

type GattungStore interface {
	reindexer
	errors.Flusher
}

type reindexer interface {
	reindexOne(sku.SkuLike) error
}
