package objekte

import (
	"code.linenisgreat.com/zit-go/src/hotel/sku"
)

type ApproximatedTyp struct {
	HasValue bool
	IsActual bool
	Typ      *sku.Transacted
}

func (a ApproximatedTyp) ActualOrNil() (actual *sku.Transacted) {
	if a.HasValue && a.IsActual {
		actual = a.Typ
	}

	return
}

func (a ApproximatedTyp) ApproximatedOrActual() *sku.Transacted {
	if !a.HasValue {
		return nil
	}

	return a.Typ
}
