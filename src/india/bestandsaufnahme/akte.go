package bestandsaufnahme

import (
	"github.com/friedenberg/zit/src/golf/sku"
)

type Akte struct {
	Skus sku.Sku2Heap
}

func MakeAkte() *Akte {
	return &Akte{
		Skus: sku.MakeSku2Heap(),
	}
}

func (a Akte) Equals(b *Akte) bool {
	if !a.Skus.Equals(b.Skus) {
		return false
	}

	return true
}

func (a *Akte) Reset(b *Akte) {
	if b == nil {
		a.Skus.Reset(nil)
	} else {
		a.Skus.Reset(&b.Skus)
	}
}
