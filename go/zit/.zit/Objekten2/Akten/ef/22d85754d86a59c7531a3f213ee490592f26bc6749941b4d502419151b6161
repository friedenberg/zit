package errors

import "golang.org/x/xerrors"

func BadRequest(err error) *errBadRequest {
	return &errBadRequest{err}
}

func BadRequestf(fmt string, args ...interface{}) *errBadRequest {
	return &errBadRequest{xerrors.Errorf(fmt, args...)}
}

func IsBadRequest(err error) bool {
	return Is(err, errBadRequest{})
}

type errBadRequest struct {
	error
}

func (e errBadRequest) ShouldShowStackTrace() bool {
	return false
}

func (e errBadRequest) Is(target error) bool {
	_, ok := target.(errBadRequest)
	return ok
}

func (e errBadRequest) Error() string {
	return e.error.Error()
}

func NewNormal(v string) errNormal {
	return errNormal{string: v}
}

type errNormal struct {
	string
}

func (e errNormal) ShouldShowStackTrace() bool {
	return false
}

func (e errNormal) Is(target error) bool {
	_, ok := target.(errNormal)
	return ok
}

func (e errNormal) Error() string {
	return e.string
}
