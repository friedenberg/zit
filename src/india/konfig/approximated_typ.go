package konfig

import (
	"github.com/friedenberg/zit/src/alfa/errors"
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

func (a ApproximatedTyp) Unwrap() *typ.Transacted {
	if !a.hasValue {
		return nil
	}

	errors.TodoP0("replace with ApproximatedOrActual")
	return &a.typ
}
