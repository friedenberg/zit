package format

import (
	"errors"
	"fmt"
)

type ErrExhaustedFuncSetStringersLine struct {
	error
	string
}

func (e ErrExhaustedFuncSetStringersLine) Error() string {
	return fmt.Sprintf("exhausted FuncSetString'ers at segment: %q", e.string)
}

func (e ErrExhaustedFuncSetStringersLine) Is(target error) (ok bool) {
	_, ok = target.(ErrExhaustedFuncSetStringersLine)
	return
}

func IsErrExhaustedFuncSetStringers(err error) bool {
	return errors.Is(err, ErrExhaustedFuncSetStringersLine{})
}

func ErrExhaustedFuncSetStringersGetString(in error) (ok bool, v string) {
	var err ErrExhaustedFuncSetStringersLine

	if errors.As(in, &err) {
		v = err.string
	}

	return
}

func ErrExhaustedFuncSetStringersWithDelim(in error, delim byte) (out error) {
	var err ErrExhaustedFuncSetStringersLine

	if errors.As(in, &err) {
		err.string = fmt.Sprintf("%s%s", err.string, string([]byte{delim}))
		in = err
	}

	out = in

	return
}
