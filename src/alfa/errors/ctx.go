package errors

import (
	"errors"
	"fmt"
)

type Ctx struct {
	Err      error
	deferred []error
}

func (e *Ctx) ClearErr() {
  e.Err = nil
  // e.deferred = make([]error, 0)
}

func (e Ctx) IsEmpty() bool {
	return e.Err == nil && len(e.deferred) == 0
}

func (d *Ctx) Defer(fs ...func() error) {
	for _, f := range fs {
		if err := f(); err != nil {
			d.deferred = append(d.deferred, err)
		}
	}
}

func (in *Ctx) Wrap() {
	var normal normalError

	if As(in.Err, &normal) {
		in.Err = normal
		return
	}

	var stack errer

	if As(in.Err, &stack) {
		ok := false
		if in.Err, ok = newStackWrapError(1); ok {
			stack.errers = append(stack.errers, in.Err)
		}

		in.Err = stack
		return
	}

	stack.errers = append(stack.errers, in.Err)
	in.Err = stack

	return
}

func (in *Ctx) Wrapf(f string, values ...interface{}) {
	err := errer{
		errers: []error{
			wrapped{
				outer: errors.New(fmt.Sprintf(f, values...)),
				inner: in.Err,
			},
		},
	}

	if st, ok := newStackWrapError(1); ok {
		err.errers = append(err.errers, st)
	}

	in.Err = err

	return
}

// func (e Ctx) Error() string {
// 	if e.Err == nil && len(e.deferred) == 0 {
// 		return ""
// 	}

// 	if len(e.deferred) == 0 {
// 		return e.Err.Error()
// 	}

// 	if e.Err == nil {
// 		return fmt.Sprintf("%s", e.deferred)
// 	}

// 	return fmt.Sprintf("multiple errors!\nerr: %s\ndeferred: %s", e.Err, e.deferred)
// }

func (e Ctx) Error() error {
	if e.IsEmpty() {
		return nil
	}

	if len(e.deferred) == 0 {
		return e.Err
	}

	if e.Err == nil {
		switch len(e.deferred) {

		default:
			return nil

		case 1:
			return e.deferred[0]
		}
	}

	//TODO return combined error

	return nil
}

// func (e Ctx) Is(target error) bool {
// 	if e.Err == nil && len(e.deferred) == 0 {
// 		return false
// 	}

// 	if len(e.deferred) == 0 && ErrorHasIsMethod(e.Err) {
// 		return e.Err.(errorWithIsMethod).Is(target)
// 	}

// 	if e.Err == nil {
// 		ok := false

// 		for _, er := range e.deferred {
// 			if !ErrorHasIsMethod(er) {
// 				continue
// 			}

// 			ok = er.(errorWithIsMethod).Is(target)

// 			if ok {
// 				return ok
// 			}
// 		}
// 	}

// 	return false
// }
