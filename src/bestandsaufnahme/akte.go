package bestandsaufnahme

import (
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/golf/sku"
)

type Akte struct {
	Skus collections.MutableValueSet[sku.Sku2, *sku.Sku2]
}

func MakeAkte() *Akte {
	return &Akte{
		Skus: collections.MakeMutableValueSet[sku.Sku2, *sku.Sku2](),
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
		//TODO-P4 make more performant
		a.Skus = collections.MakeMutableValueSet[sku.Sku2, *sku.Sku2]()
	} else {
		a.Skus.Reset(b.Skus)
	}
}
