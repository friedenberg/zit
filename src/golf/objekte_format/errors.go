package objekte_format

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type errInvalidGenericFormat string

func (err errInvalidGenericFormat) Error() string {
	return fmt.Sprintf("invalid format: %q", string(err))
}

func (err errInvalidGenericFormat) Is(target error) bool {
	_, ok := target.(errInvalidGenericFormat)
	return ok
}

var errEmptyTai = errors.New("empty tai")
