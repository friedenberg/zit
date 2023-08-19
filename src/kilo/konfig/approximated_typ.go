package konfig

import "github.com/friedenberg/zit/src/india/transacted"

type ApproximatedTyp struct {
	hasValue bool
	isActual bool
	typ      transacted.Typ
}

func (a ApproximatedTyp) HasValue() bool {
	return a.hasValue
}

func (a ApproximatedTyp) ActualOrNil() (actual *transacted.Typ) {
	if a.hasValue && a.isActual {
		actual = &a.typ
	}

	return
}

func (a ApproximatedTyp) ApproximatedOrActual() *transacted.Typ {
	if !a.hasValue {
		return nil
	}

	return &a.typ
}
