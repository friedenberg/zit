package errors

import (
	"fmt"
	"strings"
)

type Multi []error

func MakeMulti(errs ...error) *Multi {
	em := Multi{}

	for _, err := range errs {
		if err != nil {
			em.Add(err)
		}
	}

	if em.Empty() {
		return nil
	}

	return &em
}

func (e Multi) Empty() (ok bool) {
	ok = len(e) == 0
	return
}

func (e *Multi) Add(err error) {
	*e = append(*e, err)
}

func (e Multi) Is(target error) (ok bool) {
	for _, err := range e {
		if ok = Is(err, target); ok {
			return
		}
	}

	return
}

func (e Multi) Error() string {
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
