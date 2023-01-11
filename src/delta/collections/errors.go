package collections

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

var ErrStopIteration = io.EOF

func IsStopIteration(err error) bool {
	return errors.Is(err, ErrStopIteration)
}

// type ErrStopIteration struct {
// }

// func (e ErrStopIteration) Error() string {
// 	return "nil pointer"
// }

// func (e ErrStopIteration) Is(target error) (ok bool) {
// 	_, ok = target.(ErrStopIteration)
// 	return
// }

type ErrNilPointer struct {
}

func (e ErrNilPointer) Error() string {
	return "nil pointer"
}

func (e ErrNilPointer) Is(target error) (ok bool) {
	_, ok = target.(ErrNilPointer)
	return
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

type ErrDoNotRepool struct{}

func (e ErrDoNotRepool) Error() string {
	return "should not repool this element"
}

func (e ErrDoNotRepool) Is(target error) (ok bool) {
	_, ok = target.(ErrDoNotRepool)
	return
}

func IsDoNotRepool(err error) bool {
	return errors.Is(err, ErrDoNotRepool{})
}
