package store_objekten

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type GattungStore[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
	T2 gattung.Identifier[T2],
	T3 gattung.IdentifierPtr[T2],
	T4 gattung.Verzeichnisse[T],
	T5 gattung.VerzeichnissePtr[T4, T],
] interface {
	reindexer

	objekte.Store[T, T1, T2, T3, T4, T5]
	objekte.StoreLogger[T, T1, T2, T3, T4, T5]
}

type reindexer interface {
	reindexOne(*transaktion.Transaktion, sku.SkuLike) error
}
