package objekte

import (
	"github.com/friedenberg/zit/src/echo/sku"
)

//Gattung is in Sku2
type Transacted2[T any, T1 ObjektePtr[T], T2 Identifier[T2], T3 IdentifierPtr[T2]] struct {
	Objekte T
	Kennung T2
	Sku     sku.Sku2
}

func (zt Transacted2[T, T1, T2, T3]) IsNew() bool {
	return zt.Sku.Kopf == zt.Sku.Schwanz && zt.Sku.TransactionIndex == 0
}

//TODO
func (a Transacted2[T, T1, T2, T3]) Equals(
	b Transacted2[T, T1, T2, T3],
) bool {
	return false
}

func (a *Transacted2[T, T1, T2, T3]) Reset(
	b *Transacted2[T, T1, T2, T3],
) {
	a.Sku.Reset(&b.Sku)
	//TODO
	// a.Kennung.Reset(b.Kennung)
	// a.Objekte.Reset(b.Objekte)
}
