package errors

import (
	"errors"
	"io"
	"os"
)

func Is(err, target error) bool {
	e := Unwrap(err)
	return errors.Is(e, target)
}

func IsTooManyOpenFiles(err error) bool {
	e := Unwrap(err)
	return e.Error() == "too many open files"
}

func IsEOF(err error) bool {
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
