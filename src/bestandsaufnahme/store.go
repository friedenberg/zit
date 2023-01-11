package bestandsaufnahme

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/golf/sku"
)

type Store[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
] interface {
	InflateObjekteFromSku(sk sku.Sku) (T1, error)
	// objekte.Store[
	// 	Objekte,
	// 	*Objekte,
	// 	sha.Sha,
	// 	*sha.Sha,
	// 	objekte.NilVerzeichnisse[Objekte],
	// 	*objekte.NilVerzeichnisse[Objekte],
	// ]

	// objekte.StoreIdReader[
	// 	Objekte,
	// 	*Objekte,
	// 	sha.Sha,
	// 	*sha.Sha,
	// 	objekte.NilVerzeichnisse[Objekte],
	// 	*objekte.NilVerzeichnisse[Objekte],
	// ]
}

type store struct {
}
