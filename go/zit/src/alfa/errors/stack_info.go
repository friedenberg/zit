package errors

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type StackTracer interface {
	error
	ShouldShowStackTrace() bool
}

type StackInfo struct {
	Package     string
	Function    string
	Filename    string
	RelFilename string
	Line        int
}

func MakeStackInfos(depth, count int) (si []StackInfo) {
	pcs := make([]uintptr, count)

	n := runtime.Callers(depth+1, pcs)

	if n <= 0 {
		return
	}

	frames := runtime.CallersFrames(pcs)

	for {
		frame, more := frames.Next()

		si = append(si, MakeStackInfoFromFrame(frame))

		if !more {
			break
		}
	}

	return
}

func MakeStackInfoFromFrame(frame runtime.Frame) (si StackInfo) {
	si.Filename = filepath.Clean(frame.File)
	si.Line = frame.Line
	si.Function = frame.Function
	si.Package, si.Function = getPackageAndFunctionName(si.Function)

	si.RelFilename, _ = filepath.Rel(store_fs, si.Filename)

	return
}

func MakeStackInfo(skip int) (si StackInfo, ok bool) {
	var pc uintptr
	pc, _, _, ok = runtime.Caller(skip + 1)

	if !ok {
		return
	}

	frames := runtime.CallersFrames([]uintptr{pc})

	frame, _ := frames.Next()
	si = MakeStackInfoFromFrame(frame)

	return
}

func getPackageAndFunctionName(v string) (p string, f string) {
	p, f = filepath.Split(v)

	idx := strings.Index(f, ".")

	if idx == -1 {
		return
	}

	p += f[:idx]

	if len(f) > idx+1 {
		f = f[idx+1:]
	}

	return
}

func (si StackInfo) String() string {
	testPrefix := ""

	if isTest {
		testPrefix = "    "
	}

	filename := si.Filename

	if si.RelFilename != "" {
		filename = si.RelFilename
	}

	// TODO-P3 determine if si.line is ever not valid
	return fmt.Sprintf(
		"# %s\n%s%s:%d",
		si.Function,
		testPrefix,
		filename,
		si.Line,
	)
}

func (si StackInfo) Wrap(in error) (err errer) {
	se := stackWrapError{StackInfo: si}

	if As(in, &err) {
		in = se
	} else {
		in = wrapped{
			outer: se,
			inner: in,
		}
	}

	err.errers = append(err.errers, in)

	return
}

func (si StackInfo) Wrapf(in error, f string, values ...interface{}) (err errer) {
	se := stackWrapError{StackInfo: si}
	err = wrapf(se, in, f, values...)

	se, _ = newStackWrapError(1)
	err = wrapf(se, err, "")

	return
}

func (si StackInfo) Errorf(f string, values ...interface{}) (err errer) {
	e := New(fmt.Sprintf(f, values...))
	se := stackWrapError{StackInfo: si}
	err = wrapf(se, e, "")

	se, _ = newStackWrapError(1)
	err = wrapf(se, err, "")

	return
}

type stackWrapError struct {
	StackInfo
	error
}

func newStackWrapError(skip int) (err stackWrapError, ok bool) {
	var si StackInfo

	if si, ok = MakeStackInfo(skip + 1); !ok {
		return
	}

	err = stackWrapError{
		StackInfo: si,
	}

	return
}

func (se stackWrapError) Unwrap() error {
	return se.error
}

// TODO use StackInfo.String() method
func (se stackWrapError) Error() string {
	sb := &strings.Builder{}

	sb.WriteString(se.StackInfo.String())

	if se.error != nil {
		sb.WriteString(": ")
		sb.WriteString(se.error.Error())
	}

	return sb.String()
}
