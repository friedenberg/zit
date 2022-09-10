package errors

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

var cwd string
var isTest bool
var verbose bool
var maxCallDepth int

func init() {
	var err error

	if cwd, err = os.Getwd(); err != nil {
		log.Panic(err)
	}

	log.SetOutput(ioutil.Discard)
}

func SetVerbose() {
	verbose = true
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	log.Print("verbose")
}

//TODO experiment with callerframe looping
func SetCallDepth(d int) {
	maxCallDepth = d
	log.Printf("maxCallDepth: %d", maxCallDepth)
}

func SetTesting() {
	SetVerbose()
	isTest = true
	//TODO use base directory for project
	cwd = filepath.Dir(filepath.Dir(cwd))
}
