package collections

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
)

var (
	ErrNilPointer        = errors.New("nil pointer")
	ErrExists            = errors.New("exists")
	MakeErrStopIteration = iter.MakeErrStopIteration
	ErrNotFound          = errNotFound("")
)

func MakeErrNotFound(value schnittstellen.Stringer) error {
	return errors.WrapN(1, errNotFound(value.String()))
}

func MakeErrNotFoundString(s string) error {
	return errors.WrapN(1, errNotFound(s))
}

func IsErrNotFound(err error) bool {
	return errors.Is(err, errNotFound(""))
}

type errNotFound string

func (e errNotFound) Error() string {
	v := string(e)

	if v == "" {
		return "not found"
	} else {
		return fmt.Sprintf("not found: %q", v)
	}
}

func (e errNotFound) Is(target error) (ok bool) {
	_, ok = target.(errNotFound)
	return
}

type ErrEmptyKey[T any] struct {
	Element T
}

func (e ErrEmptyKey[T]) Error() string {
	return "empty key"
}

func (e ErrEmptyKey[T]) Is(target error) (ok bool) {
	_, ok = target.(ErrEmptyKey[T])
	return
}
