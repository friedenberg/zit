package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
)

type GattungStore interface {
	reindexer
	errors.Flusher
}

type reindexer interface {
	reindexOne(*transaktion.Transaktion, sku.SkuLike) error
}
