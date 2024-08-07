package errors

import (
	"errors"
	"fmt"
	"os"
)

var New = errors.New

func unwrapOnce(err error) error {
	if e, ok := err.(Unwrapper); ok {
		return e
	}

	return err
}

func Split(err error) (out []error) {
	switch e := err.(type) {
	case nil:
		return []error{}

	case *multi:
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
