package errors

import (
	"fmt"
	"strings"
)

type ErrorMulti []error

func (e ErrorMulti) Empty() (ok bool) {
	ok = len(e) == 0
	return
}

func (e *ErrorMulti) Add(err error) {
	*e = append(*e, err)
}

func (e ErrorMulti) Is(target error) (ok bool) {
	for _, err := range e {
		if ok = Is(err, target); ok {
			return
		}
	}

	return
}

func (e ErrorMulti) Error() string {
	sb := &strings.Builder{}

	sb.WriteString("# Multiple Errors")
	sb.WriteString("\n")

	for i, err := range e {
		sb.WriteString(fmt.Sprintf("Error %d:\n", i+1))
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}

	return sb.String()
}
