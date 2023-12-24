package objekte_format

import "fmt"

type errInvalidGenericFormat string

func (err errInvalidGenericFormat) Error() string {
	return fmt.Sprintf("invalid format: %q", string(err))
}
