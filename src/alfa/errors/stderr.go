package errors

import (
	"fmt"
	"os"
)

func PrintErr(a ...interface{}) {
	fmt.Fprintln(
		os.Stderr,
		a...,
	)
}

func PrintErrf(f string, a ...interface{}) {
	fmt.Fprintln(
		os.Stderr,
		fmt.Sprintf(f, a...),
	)
}
