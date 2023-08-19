package konfig

import "github.com/friedenberg/zit/src/golf/sku"

type ApproximatedTyp struct {
	hasValue bool
	isActual bool
	typ      sku.TransactedTyp
}

func (a ApproximatedTyp) HasValue() bool {
	return a.hasValue
}

func (a ApproximatedTyp) ActualOrNil() (actual *sku.TransactedTyp) {
	if a.hasValue && a.isActual {
		actual = &a.typ
	}

	return
}

func (a ApproximatedTyp) ApproximatedOrActual() *sku.TransactedTyp {
	if !a.hasValue {
		return nil
	}

	return &a.typ
}
