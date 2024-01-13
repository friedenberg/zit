package errors

import (
	"fmt"
	"strings"
)

type errer struct {
	errers []error
}

func (ers errer) Unwrap() error {
	if len(ers.errers) == 0 {
		return nil
	}

	return ers.errers[0]
}

func (e errer) Error() string {
	sb := &strings.Builder{}

	for _, e := range e.errers {
		sb.WriteString(fmt.Sprintf("%v\n", e))
	}

	return sb.String()
}
