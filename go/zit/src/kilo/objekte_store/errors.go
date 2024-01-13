package objekte_store

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
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
	ok = errors.Is(err, ErrNotFound(""))
	return
}

type ErrNotFound = collections.ErrNotFound

var ErrNotFoundEmpty = ErrNotFound("")

type VerlorenAndGefundenError interface {
	error
	AddToLostAndFound(string) (string, error)
}
