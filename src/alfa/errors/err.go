package errors

import (
	"errors"
	"fmt"
)

func WrapN(n int, in error) (err errer) {
	se, _ := newStackWrapError(1 + n)
	err = wrapf(se, in, "")
	return
}

func Wrap(in error) (err errer) {
	se, _ := newStackWrapError(1)
	err = wrapf(se, in, "")
	return
}

func Wrapf(in error, f string, values ...interface{}) (err errer) {
	se, _ := newStackWrapError(1)
	err = wrapf(se, in, f, values...)
	return
}

func Errorf(f string, values ...interface{}) (err errer) {
	e := errors.New(fmt.Sprintf(f, values...))
	se, _ := newStackWrapError(1)
	err = wrapf(se, e, "")
	return
}

func wrapf(se stackWrapError, in error, f string, values ...interface{}) (err errer) {
	// TODO-P2 case where values are present but f is ""
	if f != "" {
		se.error = errors.New(fmt.Sprintf(f, values...))
	}

	if As(in, &err) {
		in = se
	} else {
		in = wrapped{
			outer: se,
			inner: in,
		}
	}

	err.errers = append(err.errers, in)

	return
}
