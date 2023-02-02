package konfig

import (
	"github.com/friedenberg/zit/src/hotel/typ"
)

type ApproximatedTyp struct {
	hasValue bool
	isActual bool
	typ      typ.Transacted
}

func (a ApproximatedTyp) HasValue() bool {
	return a.hasValue
}

func (a ApproximatedTyp) ActualOrNil() (actual *typ.Transacted) {
	if a.hasValue && a.isActual {
		actual = &a.typ
	}

	return
}

func (a ApproximatedTyp) ApproximatedOrActual() *typ.Transacted {
	if !a.hasValue {
		return nil
	}

	return &a.typ
}
