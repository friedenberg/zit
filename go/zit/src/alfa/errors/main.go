package errors

import (
	"fmt"
	"iter"

	"golang.org/x/xerrors"
)

var New = xerrors.New

func Join(es ...error) error {
	switch {
	case len(es) == 2 && es[0] == nil && es[1] == nil:
		return nil

	case len(es) == 2 && es[0] == nil:
		return es[1]

	case len(es) == 2 && es[1] == nil:
		return es[0]

	default:
		err := MakeMulti(es...)

		if err.Empty() {
			return nil
		} else {
			return err
		}
	}
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

//go:noinline
func WrapSkip(
	skip int,
	in error,
) (err *stackWrapError) {
	if in == nil {
		return
	}

	var si StackFrame
	var ok bool

	if si, ok = MakeStackFrame(skip + 1); !ok {
		panic("failed to get stack info")
	}

	err = &stackWrapError{
		StackFrame: si,
	}

	if swe, ok := in.(*stackWrapError); ok {
		err.next = swe
	} else {
		err.error = in
	}

	return
}

const thisSkip = 1

//go:noinline
func Errorf(f string, values ...interface{}) (err error) {
	err = fmt.Errorf(f, values...)
	return
}

//go:noinline
func ErrorWithStackf(f string, values ...interface{}) (err error) {
	err = WrapSkip(thisSkip, fmt.Errorf(f, values...))
	return
}

//go:noinline
func WrapN(n int, in error) (err error) {
	err = WrapSkip(n+thisSkip, in)
	return
}

//go:noinline
func Wrap(in error) error {
	return WrapSkip(thisSkip, in)
}

//go:noinline
func Wrapf(in error, f string, values ...any) error {
	if in == nil {
		return nil
	}

	return &stackWrapError{
		StackFrame: MustStackFrame(thisSkip),
		error:      fmt.Errorf(f, values...),
		next:       WrapSkip(thisSkip, in),
	}
}

// wrap the error with stack info unless it's one of the provided `except`
// errors, in which case return nil.
//
//go:noinline
func WrapExceptAsNil(in error, except ...error) (err error) {
	if in == nil {
		return
	}

	for _, anExcept := range except {
		if Is(in, anExcept) {
			return nil
		}
	}

	err = WrapSkip(thisSkip, in)

	return
}

// wrap the error with stack info unless it's one of the provided `except`
// errors, in which case return that bare error
//
//go:noinline
func WrapExcept(in error, except ...error) (err error) {
	if in == nil {
		return
	}

	for _, e := range except {
		if Is(in, e) {
			return in
		}
	}

	err = WrapSkip(thisSkip, in)

	return
}

var errImplement = New("not implemented")

//go:noinline
func Implement() (err error) {
	return WrapSkip(1, errImplement)
}

//go:noinline
func IterWrapped[T any](err error) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		var t T
		yield(t, WrapN(1, err))
	}
}
