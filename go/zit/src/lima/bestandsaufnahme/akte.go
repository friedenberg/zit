package bestandsaufnahme

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type Akte struct {
	Skus *sku.TransactedHeap
}

func MakeAkte() *Akte {
	return &Akte{
		Skus: sku.MakeTransactedHeap(),
	}
}

func (a *Akte) GetGattung() (g schnittstellen.GattungLike) {
	g = gattung.Bestandsaufnahme

	return
}

func (a *Akte) Equals(b *Akte) bool {
	if !a.Skus.Equals(b.Skus) {
		return false
	}

	return true
}

var Resetter resetter

type resetter struct{}

func (resetter) Reset(a *Akte) {
	a.Skus.Reset()
}

func (resetter) ResetWith(a, b *Akte) {
	a.Skus.ResetWith(b.Skus)
}
