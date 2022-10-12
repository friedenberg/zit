package errors

import (
	"errors"
	"fmt"
)

func Wrap(in error) (err error) {
	return wrapf(1, in, "")
}

func Wrapf(in error, f string, values ...interface{}) (err errer) {
	return wrapf(1, in, f, values...)
}

func wrapf(skip int, in error, f string, values ...interface{}) (err errer) {
	var stack errer
	se, _ := newStackWrapError(1 + skip)

	//TODO case where values are present but f is ""
	if f != "" {
		se.error = errors.New(fmt.Sprintf(f, values...))
	}

	if As(in, &stack) {
		in = se
	} else {
		in = wrapped{
			outer: se,
			inner: in,
		}
	}

	stack.errers = append(stack.errers, in)
	err = stack

	return
}

func Errorf(f string, values ...interface{}) (err errer) {
	e := errors.New(fmt.Sprintf(f, values...))
	return wrapf(1, e, "")
}
