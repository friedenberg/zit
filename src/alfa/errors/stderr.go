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
