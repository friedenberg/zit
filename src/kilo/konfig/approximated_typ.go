package konfig

import (
	"github.com/friedenberg/zit/src/hotel/sku"
)

type ApproximatedTyp struct {
	hasValue bool
	isActual bool
	typ      *sku.Transacted2
}

func (a ApproximatedTyp) HasValue() bool {
	return a.hasValue
}

func (a ApproximatedTyp) ActualOrNil() (actual *sku.Transacted2) {
	if a.hasValue && a.isActual {
		actual = a.typ
	}

	return
}

func (a ApproximatedTyp) ApproximatedOrActual() *sku.Transacted2 {
	if !a.hasValue {
		return nil
	}

	return a.typ
}
