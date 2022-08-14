package logz

import (
	"fmt"
	"io/ioutil"
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

	log.SetOutput(ioutil.Discard)
}

func SetVerbose() {
	verbose = true
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	log.Print("verbose")
}

func SetTesting() {
	SetVerbose()
	isTest = true
	//TODO use base directory for project
	cwd = filepath.Dir(filepath.Dir(cwd))
}

func LogPrefixWithCallerDepth(d int) string {
	_, filename, line, ok := runtime.Caller(d)

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

func LogPrefix() string {
	return LogPrefixWithCallerDepth(3)
}

func CallerNonEmpty(i int, v interface{}) {
	if v != nil {
		Caller(i+1, "%s", v)
	}
}

func Caller(i int, f string, vs ...interface{}) {
	if !verbose {
		return
	}

	vs = append([]interface{}{LogPrefixWithCallerDepth(i + 2)}, vs...)
	//TODO strip trailing newline and add back
	os.Stderr.WriteString(fmt.Sprintf("%s"+f+"\n", vs...))
}

var (
	Panic  = log.Panic
	Output = log.Output
	Fatal  = log.Fatal
)

func Print(vs ...interface{}) {
	if !verbose {
		return
	}

	if len(vs) == 0 {
		os.Stderr.WriteString(fmt.Sprintln(LogPrefix()))
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
	//TODO strip trailing newline and add back
	os.Stderr.WriteString(fmt.Sprintf("%s"+f+"\n", vs...))
}

func PrintDebug(vs ...interface{}) {
	if !verbose {
		return
	}

	for _, v := range vs {
		os.Stderr.WriteString(fmt.Sprintf("%s%#v\n", LogPrefix(), v))
	}
}
