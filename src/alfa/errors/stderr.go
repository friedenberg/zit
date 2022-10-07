package errors

import (
	"fmt"
	"os"
)

func PrintErr(a ...interface{}) (err error) {
	_, err = fmt.Fprintln(
		os.Stderr,
		a...,
	)

	return
}

func PrintErrf(f string, a ...interface{}) (err error) {
	_, err = fmt.Fprintln(
		os.Stderr,
		fmt.Sprintf(f, a...),
	)

	return
}

func CallerNonEmptyErr(i int, v interface{}) {
	if v != nil {
		Caller(i+1, "%s", v)
	}
}

func CallerErr(i int, f string, vs ...interface{}) {
	st, _ := MakeStackInfo(i + 1)

	vs = append([]interface{}{st}, vs...)
	//TODO strip trailing newline and add back
	os.Stderr.WriteString(fmt.Sprintf("%s"+f+"\n", vs...))
}
