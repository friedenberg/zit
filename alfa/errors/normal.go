package errors

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
