package errors

import (
	"fmt"
	"os"
)

func PrintOut(a ...interface{}) {
	fmt.Fprintln(
		os.Stdout,
		a...,
	)
}

func PrintOutf(f string, a ...interface{}) {
	fmt.Fprintln(
		os.Stdout,
		fmt.Sprintf(f, a...),
	)
}
