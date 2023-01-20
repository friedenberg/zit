package konfig

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/hotel/typ"
)

type ApproximatedTyp struct {
	isActual bool
	typ      typ.Transacted
}

func (a ApproximatedTyp) ActualOrNil() (actual *typ.Transacted) {
	if a.isActual {
		actual = &a.typ
	}

	return
}

func (a ApproximatedTyp) ApproximatedOrActual() *typ.Transacted {
	return &a.typ
}

func (a ApproximatedTyp) Unwrap() *typ.Transacted {
	errors.TodoP0("replace with ApproximatedOrActual")
	return &a.typ
}
