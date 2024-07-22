package external_store

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type ErrUnsupportedTyp ids.Type

func (e ErrUnsupportedTyp) Is(target error) bool {
	_, ok := target.(ErrUnsupportedTyp)
	return ok
}

func (e ErrUnsupportedTyp) Error() string {
	return fmt.Sprintf("unsupported typ: %q", ids.Type(e))
}
