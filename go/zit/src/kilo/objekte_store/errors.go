package objekte_store

import (
	"fmt"
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

type VerlorenAndGefundenError interface {
	error
	AddToLostAndFound(string) (string, error)
}
