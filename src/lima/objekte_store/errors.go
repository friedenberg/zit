package objekte_store

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type ErrLockRequired struct {
	Operation string
}

func (e ErrLockRequired) Is(target error) bool {
	_, ok := target.(ErrLockRequired)
	return ok
}

func (e ErrLockRequired) Error() string {
	return fmt.Sprintf(
		"lock required for operation: %q",
		e.Operation,
	)
}

func IsNotFound(err error) (ok bool) {
	ok = errors.Is(err, ErrNotFound{})
	return
}

type ErrNotFound struct {
	Id fmt.Stringer
}

func (e ErrNotFound) Is(target error) bool {
	_, ok := target.(ErrNotFound)
	return ok
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("objekte with id '%s' not found", e.Id)
}

type VerlorenAndGefundenError interface {
	error
	AddToLostAndFound(string) (string, error)
}