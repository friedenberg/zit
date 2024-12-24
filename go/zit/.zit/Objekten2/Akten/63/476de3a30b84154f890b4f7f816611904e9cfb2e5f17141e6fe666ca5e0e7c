package errors

import (
	"io"
	log_package "log"
	"os"
)

var (
	cwd          string
	isTest       bool
	maxCallDepth int
)

func init() {
	var err error

	if cwd, err = os.Getwd(); err != nil {
		log_package.Panic(err)
	}

	log_package.SetOutput(io.Discard)
}

func SetTesting() {
	isTest = true
}
