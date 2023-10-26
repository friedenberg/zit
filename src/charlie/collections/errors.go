package collections

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
)

var (
	ErrNilPointer        = errors.New("nil pointer")
	MakeErrStopIteration = iter.MakeErrStopIteration
)

type ErrNotFound string

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("not found: %q", string(e))
}

func (e ErrNotFound) Is(target error) (ok bool) {
	_, ok = target.(ErrNotFound)
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
