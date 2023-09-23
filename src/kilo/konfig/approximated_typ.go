package konfig

import (
	"github.com/friedenberg/zit/src/hotel/sku"
)

type ApproximatedTyp struct {
	hasValue bool
	isActual bool
	typ      *sku.Transacted
}

func (a ApproximatedTyp) HasValue() bool {
	return a.hasValue
}

func (a ApproximatedTyp) ActualOrNil() (actual *sku.Transacted) {
	if a.hasValue && a.isActual {
		actual = a.typ
	}

	return
}

func (a ApproximatedTyp) ApproximatedOrActual() *sku.Transacted {
	if !a.hasValue {
		return nil
	}

	return a.typ
}

func (a ApproximatedTyp) ActualOrNilSku() (actual sku.SkuLikePtr) {
	if a.hasValue && a.isActual {
		actual = a.typ
	}

	return
}

func (a ApproximatedTyp) ApproximatedOrActualSku() sku.SkuLikePtr {
	if !a.hasValue {
		return nil
	}

	return a.typ
}
