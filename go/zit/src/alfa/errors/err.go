package errors

import (
	"fmt"
)

func wrapSkip(skip int, in error) error {
	if in == nil {
		return nil
	}

	est, ok := in.(*errorStackTrace)

	if ok {
		in = nil
	} else {
		est = &errorStackTrace{}
	}

	est.addError(1+skip, in)

	return est
}

func WrapN(n int, in error) (err error) {
	err = wrapSkip(n+1, in)
	return
}

func Wrap(in error) (err error) {
	err = wrapSkip(1, in)
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

	err = wrapSkip(1, in)

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

	err = wrapSkip(1, in)

	return
}

func Wrapf(in error, f string, values ...interface{}) (err error) {
	if in == nil {
		return nil
	}

	est, ok := in.(*errorStackTrace)

	if ok {
		in = nil
	} else {
		est = &errorStackTrace{}
	}

	est.addError(1, fmt.Errorf(f, values...))
	est.addError(1, in)

	return est
}

func Errorf(f string, values ...interface{}) (err error) {
	err = wrapSkip(1, fmt.Errorf(f, values...))
	return
}
