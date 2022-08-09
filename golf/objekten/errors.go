package objekten

import (
	"fmt"

	"github.com/friedenberg/zit/bravo/id"
)

type ErrNotFound struct {
	id.Id
}

func (e ErrNotFound) Is(target error) bool {
	_, ok := target.(ErrNotFound)
	return ok
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("objekte with id '%s' not found", e.Id)
}
