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

func (a *Akte) Reset() {
	a.Skus.Reset()
}

func (a *Akte) ResetWith(b Akte) {
	a.Skus.ResetWith(b.Skus)
}
