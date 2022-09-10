package errors

import (
	"fmt"
	"log"
	"os"
)

func CallerNonEmpty(i int, v interface{}) {
	if v != nil {
		Caller(i+1, "%s", v)
	}
}

func Caller(i int, f string, vs ...interface{}) {
	if !verbose {
		return
	}

	st, _ := MakeStackInfo(i + 1)

	vs = append([]interface{}{st}, vs...)
	//TODO strip trailing newline and add back
	os.Stderr.WriteString(fmt.Sprintf("%s"+f+"\n", vs...))
}

var (
	//TODO add native methods
	Panic  = log.Panic
	Output = log.Output
	Fatal  = log.Fatal
)

func Print(vs ...interface{}) {
	if !verbose {
		return
	}

	si, _ := MakeStackInfo(1)

	if len(vs) == 0 {
		os.Stderr.WriteString(fmt.Sprintln(si))
	}

	for _, v := range vs {
		os.Stderr.WriteString(fmt.Sprintln(si, v))
	}
}

func Printf(f string, vs ...interface{}) {
	if !verbose {
		return
	}

	si, _ := MakeStackInfo(1)

	vs = append([]interface{}{si}, vs...)
	//TODO strip trailing newline and add back
	os.Stderr.WriteString(fmt.Sprintf("%s"+f+"\n", vs...))
}

func PrintDebug(vs ...interface{}) {
	if !verbose {
		return
	}

	si, _ := MakeStackInfo(1)

	for _, v := range vs {
		os.Stderr.WriteString(fmt.Sprintf("%s%#v\n", si, v))
	}
}
