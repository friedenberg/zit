package errors

import "fmt"

type wrapError struct {
	err    error
	frames []_Frame
}

var Errorf = _Errorf

func Error(err error) error {
	var normal StackTracer

	if As(err, &normal) {
		return err
	}

	var wrapped wrapError

	if As(err, &wrapped) {
		return err
	}

	return &wrapError{
		err:    err,
		frames: []_Frame{_Caller(2)},
	}
}

func (e wrapError) Error() string {
	return e.err.Error()
}

func (e wrapError) Format(s fmt.State, v rune) {
	_FormatError(e, s, v)
}

func (e wrapError) FormatError(p _Printer) (next error) {
	for _, f := range e.frames {
		f.Format(p)
	}

	return e.err
}

func (e wrapError) Unwrap() error {
	return e.err
}
