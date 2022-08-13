package errors

import (
	"errors"
	"io"
	"os"

	"golang.org/x/xerrors"
)

type unwrappable interface {
	Unwrap() error
}

func As(err error, target interface{}) bool {
	return xerrors.As(err, target)
}

func Is(err, target error) bool {
	e := Unwrap(err)
	return errors.Is(e, target)
	// ok := err.(target)
	// // ok := xerrors.Is(err, target)

	// if ok {
	// 	return true
	// }

	// if e, ok := err.(unwrappable); ok {
	// 	return Is(e.Unwrap(), target)
	// }

	// return false
}

func IsEOF(err error) bool {
	e := Unwrap(err)
	return Is(e, io.EOF)
}

func IsNotExist(err error) bool {
	e := Unwrap(err)
	return os.IsNotExist(e)
}

func Unwrap(err error) error {
	if e, ok := err.(unwrappable); ok {
		return Unwrap(e.Unwrap())
	}

	return err
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
