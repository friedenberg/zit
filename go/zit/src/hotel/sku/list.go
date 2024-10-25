package sku

type List = TransactedHeap

func MakeList() *List {
	return MakeTransactedHeap()
}

var ResetterList resetterList

type resetterList struct{}

func (resetterList) Reset(a *List) {
	a.Reset()
}

func (resetterList) ResetWith(a, b *List) {
	a.ResetWith(b)
}
