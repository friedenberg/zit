package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var cwd string
var isTest bool
var verbose bool

func init() {
	var err error

	if cwd, err = os.Getwd(); err != nil {
		log.Panic(err)
	}
}

func SetVerbose() {
  verbose = true
}

func SetTesting() {
  SetVerbose()
	isTest = true
	//TODO use base directory for project
	cwd = filepath.Dir(filepath.Dir(cwd))
}

func LogPrefix() string {
	_, filename, line, ok := runtime.Caller(2)

	if !ok {
		return ""
	}

	filename = filepath.Clean(filename)

	var p string
	var err error

	if p, err = filepath.Rel(cwd, filename); err != nil {
		return ""
	}

	testPrefix := ""

	if isTest {
		testPrefix = "    "
	}

	return fmt.Sprintf("%s%s:%d: ", testPrefix, p, line)
}

func Print(vs ...interface{}) {
  if !verbose {
    return 
  }

	for _, v := range vs {
		os.Stderr.WriteString(fmt.Sprintln(LogPrefix(), v))
	}
}

func Printf(f string, vs ...interface{}) {
  if !verbose {
    return 
  }

	vs = append([]interface{}{LogPrefix()}, vs...)
	os.Stderr.WriteString(fmt.Sprintf("%s"+f, vs...))
}

func PrintDebug(vs ...interface{}) {
  if !verbose {
    return 
  }

	for _, v := range vs {
		os.Stderr.WriteString(fmt.Sprintf("%s%#v\n", LogPrefix(), v))
	}
}
