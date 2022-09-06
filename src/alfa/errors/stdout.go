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
