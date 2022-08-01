package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type errer struct {
	errers []error
}

func (ers errer) Unwrap() error {
	if len(ers.errers) == 0 {
		return nil
	}

	return ers.errers[0]
}

type wrapped struct {
	outer, inner error
}

func (e wrapped) Error() string {
	return fmt.Sprintf("%s: %s", e.outer, e.inner)
}

func (e wrapped) Unwrap() error {
	return e.inner
}

func Wrapped(in error, f string, values ...interface{}) (err errer) {
	err = errer{
		errers: []error{
			wrapped{
				outer: errors.New(fmt.Sprintf(f, values...)),
				inner: in,
			},
		},
	}

	if st, ok := newStackWrapError(1); ok {
		err.errers = append(err.errers, st)
	}

	return
}

func Errorf(f string, values ...interface{}) (err errer) {
	err = errer{
		errers: []error{
			errors.New(fmt.Sprintf(f, values...)),
		},
	}

	if st, ok := newStackWrapError(1); ok {
		err.errers = append(err.errers, st)
	}

	return
}

// func WithInfo(in error, info string) (out error) {

// }
func newStackWrapError(skip int) (err stackWrapError, ok bool) {
	var (
		pc   uintptr
		file string
		line int
	)

	pc, file, line, ok = runtime.Caller(skip + 1)

	if !ok {
		return
	}

	frames := runtime.CallersFrames([]uintptr{pc})

	frame, _ := frames.Next()

	err = stackWrapError{
		Frame: frame,
		file:  file,
		line:  line,
	}

	return
}

type stackWrapError struct {
	runtime.Frame
	file string
	line int
}

func (se stackWrapError) Error() string {
	return fmt.Sprintf("- %s\n  %s:%d", se.Frame.Function, se.file, se.line)
}

func Error(in error) (err error) {
	var normal normalError

	if As(in, &normal) {
		err = normal
		return
	}

	var stack errer

	if As(in, &stack) {
		if in, ok := newStackWrapError(1); ok {
			stack.errers = append(stack.errers, in)
		}

		err = stack
		return
	}

	stack.errers = append(stack.errers, in)
	err = stack

	return
}

func (e errer) Error() string {
	sb := &strings.Builder{}

	for _, e := range e.errers {
		sb.WriteString(fmt.Sprintf("%v\n", e))
	}

	return sb.String()
}

// func (e errer) Format(s fmt.State, v rune) {
// 	xerrors.FormatError(e, s, v)
// }

// func (e errer) FormatError(p xerrors.Printer) (next error) {
// 	for _, f := range e.frames {
// 		f.Format(p)
// 	}

// 	return e.err
// }

// func (e errer) Unwrap() error {
// 	return e.err
// }
