package collections

import (
	"github.com/friedenberg/zit/src/alfa/errors"
)

var (
	errStopIteration = errors.New("stop iteration")
	ErrNilPointer    = errors.New("nil pointer")
	ErrDoNotRepool   = errors.New("do not repool")
)

func MakeErrStopIteration() error {
	if errors.IsVerbose() {
		return errors.WrapN(2, errStopIteration)
	} else {
		return errStopIteration
	}
}

func IsStopIteration(err error) bool {
	if errors.Is(err, errStopIteration) {
		errors.Log().Printf("stopped iteration at %s", err)
		return true
	}

	return false
}

func IsDoNotRepool(err error) bool {
	return errors.Is(err, ErrDoNotRepool)
}

type ErrNotFound struct {
}

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
