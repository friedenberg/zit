package errors

import (
	"fmt"
	"strings"
)

type errorStackTrace struct {
	errors []stackWrapError
}

func (ers *errorStackTrace) add(e stackWrapError) {
	ers.errors = append(ers.errors, e)
}

func (ers *errorStackTrace) addError(skip int, e error) {
	err, _ := newStackWrapError(skip+1, e)
	ers.errors = append(ers.errors, err)
}

func (ers errorStackTrace) Unwrap() error {
	if len(ers.errors) == 0 {
		return nil
	}

	return ers.errors[0]
}

func (e errorStackTrace) Error() string {
	sb := &strings.Builder{}

	for i := len(e.errors) - 1; i >= 0; i-- {
		e := e.errors[i]
		fmt.Fprintf(sb, "%v\n", e)
	}

	return sb.String()
}
