package errors

import (
	"fmt"
	"os"
)

func PrintOut(a ...interface{}) (err error) {
	_, err = fmt.Fprintln(
		os.Stdout,
		a...,
	)

	return
}

func PrintOutf(f string, a ...interface{}) (err error) {
	_, err = fmt.Fprintln(
		os.Stdout,
		fmt.Sprintf(f, a...),
	)

	return
}

func CallerNonEmptyOut(i int, v interface{}) {
	if v != nil {
		Caller(i+1, "%s", v)
	}
}

func CallerOut(i int, f string, vs ...interface{}) {
	st, _ := MakeStackInfo(i + 1)

	vs = append([]interface{}{st}, vs...)
	//TODO strip trailing newline and add back
	os.Stderr.WriteString(fmt.Sprintf("%s"+f+"\n", vs...))
}
