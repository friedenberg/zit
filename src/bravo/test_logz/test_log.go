package test_logz

import (
	"fmt"
	"os"
	"testing"

	"github.com/friedenberg/zit/src/alfa/errors"
)

var (
	Print         = errors.Log().Print
	Printf        = errors.Log().Printf
	MakeStackInfo = errors.MakeStackInfo
)

type (
	StackInfo = errors.StackInfo
)

type T struct {
	*testing.T
	skip int
}

func (t T) Skip(skip int) T {
	return T{
		T:    t.T,
		skip: t.skip + skip,
	}
}

func (t T) logf(skip int, format string, args ...interface{}) {
	errors.SetTesting()
	si, _ := MakeStackInfo(t.skip + 1 + skip)
	args = append([]interface{}{si}, args...)
	os.Stderr.WriteString(fmt.Sprintf("%s"+format+"\n", args...))
}

func (t T) errorf(skip int, format string, args ...interface{}) {
	t.logf(skip+1, format, args...)
	t.Fail()
}

func (t T) fatalf(skip int, format string, args ...interface{}) {
	errors.SetTesting()
	si, _ := MakeStackInfo(t.skip + 1 + skip)
	args = append([]interface{}{si}, args...)
	os.Stderr.WriteString(fmt.Sprintf("%s"+format+"\n", args...))
	t.FailNow()
}

func (t T) Logf(format string, args ...interface{}) {
	t.logf(1, format, args...)
}

func (t T) Errorf(format string, args ...interface{}) {
	t.errorf(1, format, args...)
}

func (t T) Fatalf(format string, args ...interface{}) {
	t.fatalf(1, format, args...)
}

// TODO-P3 move to AssertNotEqual
func (t T) NotEqual(a, b any) {
	format := "\nexpected: %v\n  actual: %v"
	t.errorf(1, format, a, b)
}

func (t T) AssertNoError(err error) {
	t.T.Helper()

	if err != nil {
		t.fatalf(1, "expected no error but got %q", err)
	}
}

func (t T) AssertEOF(err error) {
	t.T.Helper()

	if !errors.IsEOF(err) {
		t.fatalf(1, "expected EOF but got none")
	}
}

func (t T) AssertError(err error) {
	t.T.Helper()

	if err == nil {
		t.fatalf(1, "expected an error but got none")
	}
}
