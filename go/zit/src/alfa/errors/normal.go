package errors

import "golang.org/x/xerrors"

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

func BadRequest(err error) *errBadRequest {
	return &errBadRequest{err}
}

func BadRequestf(fmt string, args ...interface{}) *errBadRequest {
	return &errBadRequest{xerrors.Errorf(fmt, args...)}
}

func IsBadRequest(err error) bool {
	return Is(err, errBadRequest{})
}
