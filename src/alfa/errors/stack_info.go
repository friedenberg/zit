package errors

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type StackInfo struct {
	filename    string
	relFilename string
	line        int
	depth       int
	pc          uintptr
}

func MakeStackInfo(d int) (si StackInfo, ok bool) {
	si.depth = d
	si.pc, si.filename, si.line, ok = runtime.Caller(d + 1)

	if ok {
		si.filename = filepath.Clean(si.filename)
	}

	// var err error

	// if si.relFilename, err = filepath.Rel(cwd, si.filename); err != nil {
	// 	ok = false
	// 	return
	// }

	return
}

func (si StackInfo) String() string {
	testPrefix := ""

	if isTest {
		testPrefix = "    "
	}

	filename := si.filename

	if si.relFilename != "" {
		filename = si.relFilename
	}

	return fmt.Sprintf("%s%s:%d: ", testPrefix, filename, si.line)
}
