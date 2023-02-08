package errors

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/xerrors"
)

var New = xerrors.New

func As(err error, target interface{}) bool {
	return xerrors.As(err, target)
}

func unwrapOnce(err error) error {
	if e, ok := err.(Unwrapper); ok {
		return e
	}

	return err
}

func Unwrap(err error) error {
	if e, ok := err.(Unwrapper); ok {
		return Unwrap(e.Unwrap())
	}

	return err
}

func Split(err error) (out []error) {
	switch e := err.(type) {
	case nil:
		return []error{}

	case *multi:
		out = e.Errors()

	case multi:
		out = e.Errors()

	default:
		out = []error{err}
	}

	return
}

func Join(es ...error) error {
	switch {
	case len(es) == 2 && es[0] == nil && es[1] == nil:
		return nil

	case len(es) == 2 && es[0] == nil:
		return es[1]

	case len(es) == 2 && es[1] == nil:
		return es[0]

	default:
    err := MakeMulti(es...)

    if err.Empty() {
      return nil
    } else {
      return err
    }
	}
}

func DeferredFlusher(
	err *error,
	f Flusher,
) {
	if err1 := f.Flush(); err1 != nil {
		if err == nil {
			panic(err)
		} else {
			*err = Join(*err, err1)
		}
	}
}

func DeferredCloser(
	err *error,
	c io.Closer,
) {
	if err1 := c.Close(); err1 != nil {
		if err == nil {
			panic(err)
		} else {
			*err = Join(*err, err1)
		}
	}
}

func Deferred(
	err *error,
	ef func() error,
) {
	if err1 := ef(); err1 != nil {
		if err == nil {
			panic(err)
		} else {
			*err = Join(*err, err1)
		}
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
		*err = Join(*err, err1)
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
