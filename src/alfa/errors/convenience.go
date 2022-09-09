package errors

import (
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/xerrors"
)

type errorWithIsMethod interface {
	error
	Is(error) bool
}

func ErrorHasIsMethod(err error) bool {
	_, ok := err.(errorWithIsMethod)

	return ok
}

type unwrappable interface {
	error
	Unwrap() error
}

func As(err error, target interface{}) bool {
	return xerrors.As(err, target)
}

func Is(err, target error) bool {
	e := Unwrap(err)
	return errors.Is(e, target)
}

func IsEOF(err error) bool {
	return Is(err, io.EOF)
}

func IsNotExist(err error) bool {
	e := Unwrap(err)
	return os.IsNotExist(e)
}

func unwrapOnce(err error) error {
	if e, ok := err.(unwrappable); ok {
		return e
	}

	return err
}

func Unwrap(err error) error {
	if e, ok := err.(unwrappable); ok {
		return Unwrap(e.Unwrap())
	}

	return err
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

func Err(a ...interface{}) {
	fmt.Fprintln(
		os.Stderr,
		a...,
	)
}
