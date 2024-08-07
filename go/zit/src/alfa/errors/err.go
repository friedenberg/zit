package errors

import (
	"fmt"
)

func wrapSkip(skip int, in error) error {
	if in == nil {
		return nil
	}

	var out *stackWrapError

	if swe, ok := in.(*stackWrapError); ok {
		out = newStackWrapError(skip+1, nil, swe)
	} else {
		out = newStackWrapError(skip+1, in, nil)
	}

	return out
}

const thisSkip = 1

func WrapN(n int, in error) (err error) {
	err = wrapSkip(n+thisSkip, in)
	return
}

func Wrap(in error) (err error) {
	err = wrapSkip(thisSkip, in)
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

	err = wrapSkip(thisSkip, in)

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

	err = wrapSkip(thisSkip, in)

	return
}

func Wrapf(in error, f string, values ...interface{}) error {
	if in == nil {
		return nil
	}

	var inner *stackWrapError

	if swe, ok := in.(*stackWrapError); ok {
		inner = newStackWrapError(1, nil, swe)
	} else {
		inner = newStackWrapError(1, in, nil)
	}

	return newStackWrapError(1, fmt.Errorf(f, values...), inner)
}

func Errorf(f string, values ...interface{}) (err error) {
	err = wrapSkip(thisSkip, fmt.Errorf(f, values...))
	return
}
