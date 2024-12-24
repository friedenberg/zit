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

func MakeStackInfoFromFrame(frame runtime.Frame) (si StackInfo) {
	si.Filename = filepath.Clean(frame.File)
	si.Line = frame.Line
	si.Function = frame.Function
	si.Package, si.Function = getPackageAndFunctionName(si.Function)

	si.RelFilename, _ = filepath.Rel(cwd, si.Filename)

	return
}

func MustStackInfo(skip int) StackInfo {
	si, ok := MakeStackInfo(skip + 1)

	if !ok {
		panic("stack unavailable")
	}

	return si
}

func MakeStackInfo(skip int) (si StackInfo, ok bool) {
	var pc uintptr
	pc, _, _, ok = runtime.Caller(skip + 1) // 0 is self

	if !ok {
		return
	}

	frames := runtime.CallersFrames([]uintptr{pc})

	frame, _ := frames.Next()
	si = MakeStackInfoFromFrame(frame)

	if si.Function == "Wrap" {
		panic(fmt.Sprintf("Parent Wrap included in stack. Skip: %d", skip))
	}

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

func (si StackInfo) StringNoFunctionName() string {
	filename := si.Filename

	if si.RelFilename != "" {
		filename = si.RelFilename
	}

	return fmt.Sprintf(
		"%s|%d|",
		filename,
		si.Line,
	)
}

func (si StackInfo) Wrap(in error) (err error) {
	return &stackWrapError{
		StackInfo: si,
		error:     in,
	}
}

func (si StackInfo) Wrapf(in error, f string, values ...interface{}) (err error) {
	return &stackWrapError{
		StackInfo: si,
		error:     fmt.Errorf(f, values...),
		next: &stackWrapError{
			StackInfo: si,
			error:     in,
		},
	}
}

func (si StackInfo) Errorf(f string, values ...interface{}) (err error) {
	return &stackWrapError{
		StackInfo: si,
		error:     fmt.Errorf(f, values...),
	}
}

type stackWrapError struct {
	StackInfo
	error

	next *stackWrapError
}

func (se *stackWrapError) Unwrap() error {
	if se.next == nil {
		return se.error
	} else {
		return se.next.Unwrap()
	}
}

func (se *stackWrapError) UnwrapAll() []error {
	switch {
	case se.next != nil && se.error != nil:
		return []error{se.error, se.next}

	case se.next != nil:
		return []error{se.next}

	case se.error != nil:
		return []error{se.error}

	default:
		return nil
	}
}

func (se *stackWrapError) writeError(sb *strings.Builder) {
	sb.WriteString(se.StackInfo.String())

	if se.error != nil {
		sb.WriteString(": ")
		sb.WriteString(se.error.Error())
	}

	if se.next != nil {
		sb.WriteString("\n")
		se.next.writeError(sb)
	}

	if se.next == nil && se.error == nil {
		sb.WriteString("zit/alfa/errors/stackWrapError: both next and error are nil.")
		sb.WriteString("zit/alfa/errors/stackWrapError: this usually means that some nil error was wrapped in the error stack.")
	}
}

func (se stackWrapError) Error() string {
	sb := &strings.Builder{}
	se.writeError(sb)
	return sb.String()
}
