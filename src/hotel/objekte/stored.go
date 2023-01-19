package objekte

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/golf/sku"
)

type Stored[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
	T2 schnittstellen.Id[T2],
	T3 schnittstellen.IdPtr[T2],
] struct {
	Objekte T
	Sku     sku.External[T2, T3]
}
