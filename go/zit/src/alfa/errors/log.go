package errors

import (
	"io"
	log_package "log"
	"os"
)

type Logger interface {
	// Fatal(v ...interface{})
	// Fatalf(format string, v ...interface{})

	// Panic(v ...interface{})
	// Panicf(format string, v ...interface{})
	// Panicln(v ...interface{})

	Print(v ...interface{})
	Printf(format string, v ...interface{})

	// Output(calldepth int, s string) error

	// Prefix() string
	// SetPrefix(prefix string)
}

var (
	cwd          string
	isTest       bool
	verbose      bool
	maxCallDepth int
)

func init() {
	var err error

	if cwd, err = os.Getwd(); err != nil {
		log_package.Panic(err)
	}

	log_package.SetOutput(io.Discard)
}

func SetVerbose() {
	printerLog.on = true
	printerDebug.on = true
	verbose = true
	log_package.SetFlags(log_package.LstdFlags | log_package.Lshortfile | log_package.Lmicroseconds)
	log_package.Print("verbose")
}

func IsVerbose() bool {
	return verbose
}

func SetTesting() {
	SetVerbose()
	isTest = true
}
