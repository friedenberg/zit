package inventory_list

import (
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type InventoryList = sku.TransactedHeap

func MakeInventoryList() *InventoryList {
	return sku.MakeTransactedHeap()
}

var Resetter resetter

type resetter struct{}

func (resetter) Reset(a *InventoryList) {
	a.Reset()
}

func (resetter) ResetWith(a, b *InventoryList) {
	a.ResetWith(b)
}
