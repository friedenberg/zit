package collections

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/iter"
)

var (
	ErrNilPointer        = errors.New("nil pointer")
	ErrDoNotRepool       = errors.New("do not repool")
	MakeErrStopIteration = iter.MakeErrStopIteration
	IsStopIteration      = iter.IsStopIteration
)

func IsDoNotRepool(err error) bool {
	return errors.Is(err, ErrDoNotRepool)
}

type ErrNotFound struct{}

func (e ErrNotFound) Error() string {
	return "not found"
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
