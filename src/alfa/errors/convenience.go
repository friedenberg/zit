package errors

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/xerrors"
)

type errorWithIsMethod interface {
	error
	Is(error) bool
}

func ErrorHasIsMethod(err error) bool {
	_, ok := err.(errorWithIsMethod)

	return ok
}

type unwrappable interface {
	error
	Unwrap() error
}

func As(err error, target interface{}) bool {
	return xerrors.As(err, target)
}

func unwrapOnce(err error) error {
	if e, ok := err.(unwrappable); ok {
		return e
	}

	return err
}

func Unwrap(err error) error {
	if e, ok := err.(unwrappable); ok {
		return Unwrap(e.Unwrap())
	}

	return err
}

func DeferredCloser(
	err *error,
	c io.Closer,
) {
	if err1 := c.Close(); err1 != nil {
		*err = MakeMulti(*err, err1)
	}
}

func Deferred(
	err *error,
	ef func() error,
) {
	if err1 := ef(); err1 != nil {
		*err = MakeMulti(*err, err1)
	}
}

func DeferredChanError(
	err *error,
	ch <-chan error,
) {
	var err1 error

	select {
	case err1 = <-ch:
	}

	if err1 != nil {
		*err = MakeMulti(*err, err1)
	}
}

func DeferredChan(
	ch chan<- error,
	f func() error,
) {
	if err := f(); err != nil {
		ch <- err
	}
}

func PanicIfError(err interface{}) {
	if err == nil {
		return
	}

	switch t := err.(type) {
	case func() error:
		PanicIfError(t())
	case error:
		panic(t)
	}
}

func Error(err error) {
	if err == nil {
		return
	}

	fmt.Fprintln(
		os.Stderr,
		fmt.Sprintf("%+v", err),
	)
}

func Errf(f string, a ...interface{}) {
	fmt.Fprintln(
		os.Stderr,
		fmt.Sprintf(f, a...),
	)
}
