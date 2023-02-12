package errors

import (
	"fmt"
	"strings"
	"sync"
)

type Multi interface {
	Add(error)
	Empty() bool
}

type multi struct {
	lock    sync.Locker
	chOnErr chan struct{}
	slice   []error
}

func MakeMulti(errs ...error) (em *multi) {
	em = &multi{
		lock:    &sync.Mutex{},
		chOnErr: make(chan struct{}),
		slice:   make([]error, 0, len(errs)),
	}

	for _, err := range errs {
		if err != nil {
			em.Add(err)
		}
	}

	return
}

// TODO-P4 determine why this didn't work
// func (e *multi) Combine(
// 	err *error,
// ) {
// 	if !e.Empty() && *err != nil {
// 		e.Add(*err)
// 		*err = e
// 	}
// }

func (e multi) ChanOnErr() <-chan struct{} {
	return e.chOnErr
}

func (e multi) Len() int {
	e.lock.Lock()
	defer e.lock.Unlock()

	return len(e.slice)
}

func (e multi) Empty() (ok bool) {
	ok = e.Len() == 0
	return
}

func (e *multi) merge(err multi) {
	e.lock.Lock()

	l := len(e.slice)

	e.slice = append(e.slice, err.slice...)

	if len(e.slice) > l && l == 0 {
		close(e.chOnErr)
	}

	e.lock.Unlock()
}

func (e *multi) Add(err error) {
	// if err == nil {
	// 	panic("trying to add nil error")
	// }

	if err != nil {
		if e == nil {
			// panic("trying to add to nil multi error")
			e = MakeMulti(err)
			return
		}

		switch e1 := Unwrap(err).(type) {
		case multi:
			e.merge(e1)

		case *multi:
			e.merge(*e1)

		default:
			e.lock.Lock()

			l := len(e.slice)

			e.slice = append(e.slice, err)

			if len(e.slice) > l && l == 0 {
				close(e.chOnErr)
			}

			e.lock.Unlock()
		}
	}
}

func (e multi) Is(target error) (ok bool) {
	for _, err := range e.Errors() {
		if ok = Is(err, target); ok {
			return
		}
	}

	return
}

func (e multi) Errors() (out []error) {
	e.lock.Lock()
	defer e.lock.Unlock()
	copy(out, e.slice)
	return
}

func (e multi) Error() string {
	sb := &strings.Builder{}

	sb.WriteString("# Multiple Errors")
	sb.WriteString("\n")

	e.lock.Lock()
	defer e.lock.Unlock()

	for i, err := range e.slice {
		sb.WriteString(fmt.Sprintf("Error %d:\n", i+1))
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}

	return sb.String()
}
