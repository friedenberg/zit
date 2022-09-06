package errors

import "fmt"

func Defer(d *Deferred, f func() error) {
	d.deferred = f()
}

type Deferred struct {
	err      error
	deferred error
}

func (e Deferred) Error() string {
	if e.err == nil && e.deferred == nil {
		return ""
	}

	if e.deferred == nil {
		return e.err.Error()
	}

	if e.err == nil {
		return e.deferred.Error()
	}

	return fmt.Sprintf("multiple errors!\nerr: %s\ndeferred: %s", e.err, e.deferred)
}

func (e Deferred) Unwrap() error {
	if e.err == nil && e.deferred == nil {
		return nil
	}

	if e.deferred == nil {
		return e.err
	}

	if e.err == nil {
		return e.deferred
	}

	return nil
}

func (e Deferred) Is(target error) bool {
	if e.err == nil && e.deferred == nil {
		return false
	}

	if e.deferred == nil && ErrorHasIsMethod(e.err) {
		return e.err.(errorWithIsMethod).Is(target)
	}

	if e.err == nil && ErrorHasIsMethod(e.deferred) {
		return e.deferred.(errorWithIsMethod).Is(target)
	}

	return false
}
