package errors

import (
	"fmt"
)

func WrapN(n int, in error) (err error) {
	if in == nil {
		return
	}

	var est errorStackTrace

	if As(in, &est) {
		in = nil
	}

	est.addError(2+n, in)
	err = &est

	return
}

func Wrap(in error) (err error) {
	if in == nil {
		return
	}

	var est errorStackTrace

	if As(in, &est) {
		in = nil
	}

	est.addError(2, in)
	err = &est

	return
}

func WrapExceptAsNil(in error, except ...error) (err error) {
	if in == nil {
		return
	}

	for _, e := range except {
		if in == e {
			return nil
		}
	}

	var est errorStackTrace

	if As(in, &est) {
		in = nil
	}

	est.addError(2, in)
	err = est

	return
}

func WrapExcept(in error, except ...error) (err error) {
	if in == nil {
		return
	}

	for _, e := range except {
		if in == e {
			return in
		}
	}

	var est errorStackTrace

	if As(in, &est) {
		in = nil
	}

	est.addError(2, in)
	err = est

	return
}

func Wrapf(in error, f string, values ...interface{}) (est errorStackTrace) {
	if As(in, &est) {
		in = nil
	}

	est.addError(1, in)
	est.addError(1, fmt.Errorf(f, values...))

	return
}

func Errorf(f string, values ...interface{}) (est errorStackTrace) {
	est.addError(2, fmt.Errorf(f, values...))
	return
}
