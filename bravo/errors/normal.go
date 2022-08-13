package errors

import "golang.org/x/xerrors"

type StackTracer interface {
	error
	ShouldShowStackTrace() bool
}

type normalError struct {
	error
}

func (e normalError) ShouldShowStackTrace() bool {
	return false
}

func (e normalError) Error() string {
	return e.error.Error()
}

func Normal(err error) *normalError {
	return &normalError{err}
}

func Normalf(fmt string, args ...interface{}) *normalError {
	return &normalError{xerrors.Errorf(fmt, args...)}
}
