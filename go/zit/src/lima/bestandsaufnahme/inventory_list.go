package bestandsaufnahme

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type InventoryList struct {
	Skus *sku.TransactedHeap
}

func MakeInventoryList() *InventoryList {
	return &InventoryList{
		Skus: sku.MakeTransactedHeap(),
	}
}

func (a *InventoryList) GetGattung() (g interfaces.Genre) {
	g = genres.InventoryList

	return
}

func (a *InventoryList) Equals(b *InventoryList) bool {
	if !a.Skus.Equals(b.Skus) {
		return false
	}

	return true
}

var Resetter resetter

type resetter struct{}

func (resetter) Reset(a *InventoryList) {
	a.Skus.Reset()
}

func (resetter) ResetWith(a, b *InventoryList) {
	a.Skus.ResetWith(b.Skus)
}
