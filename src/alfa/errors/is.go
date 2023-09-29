package errors

import (
	"errors"
	"io"
	"os"
	"syscall"
)

func As(err error, target any) bool {
	es := Split(err)

	switch len(es) {
	case 0:
		return false

	case 1:
		return errors.As(Unwrap(es[0]), target)

	default:
		for _, e := range es {
			if As(e, target) {
				return true
			}
		}
	}

	return false
}

func Is(err, target error) bool {
	es := Split(err)

	switch len(es) {
	case 0:
		return false

	case 1:
		return errors.Is(Unwrap(es[0]), target)

	default:
		for _, e := range es {
			if Is(e, target) {
				return true
			}
		}
	}

	return false
}

func IsErrno(err error, target syscall.Errno) (ok bool) {
	var errno syscall.Errno

	if !As(err, &errno) {
		return
	}

	ok = errno == target

	return
}

func IsBrokenPipe(err error) bool {
	return IsErrno(err, syscall.EPIPE)
}

func IsTooManyOpenFiles(err error) bool {
	e := Unwrap(err)
	return e.Error() == "too many open files"
}

func IsEOF(err error) bool {
	if err == nil {
		return false
	}

	return Is(err, io.EOF)
}

func IsExist(err error) bool {
	e := Unwrap(err)
	return os.IsExist(e)
}

func IsNotExist(err error) bool {
	e := Unwrap(err)
	return os.IsNotExist(e)
}

func IsAsNilOrWrapf(
	err error,
	target error,
	format string,
	values ...interface{},
) (out error) {
	if Is(err, target) {
		out = nil
	} else {
		out = Wrapf(err, format, values...)
	}

	return
}
