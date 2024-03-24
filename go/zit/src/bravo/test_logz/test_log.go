package test_logz

import (
	"fmt"
	"os"
	"testing"

	"code.linenisgreat.com/zit/src/alfa/errors"
)

var (
	Print  = errors.Log().Print
	Printf = errors.Log().Printf
)

type T struct {
	*testing.T
	skip int
}

func (t *T) SkipTest(args ...any) {
	t.log(1, args...)
	t.SkipNow()
}

func (t *T) Skip(skip int) T {
	return T{
		T:    t.T,
		skip: t.skip + skip,
	}
}

func (t *T) log(skip int, args ...interface{}) {
	errors.SetTesting()
	si := MakeStackInfo(t, t.skip+1+skip)
	args = append([]interface{}{si}, args...)
	fmt.Fprintln(os.Stderr, args...)
}

func (t *T) logf(skip int, format string, args ...interface{}) {
	errors.SetTesting()
	si := MakeStackInfo(t, t.skip+1+skip)
	args = append([]interface{}{si}, args...)
	fmt.Fprintf(os.Stderr, "%s"+format+"\n", args...)
}

func (t *T) errorf(skip int, format string, args ...interface{}) {
	t.logf(skip+1, format, args...)
	t.Fail()
}

func (t *T) fatalf(skip int, format string, args ...interface{}) {
	t.logf(skip+1, format, args...)
	t.FailNow()
}

func (t *T) Log(args ...interface{}) {
	t.log(1, args...)
}

func (t *T) Logf(format string, args ...interface{}) {
	t.logf(1, format, args...)
}

func (t *T) Errorf(format string, args ...interface{}) {
	t.Helper()
	t.errorf(1, format, args...)
}

func (t *T) Fatalf(format string, args ...interface{}) {
	t.Helper()
	t.fatalf(1, format, args...)
}

// TODO-P3 move to AssertNotEqual
func (t *T) NotEqual(a, b any) {
	format := "\nexpected: %q\n  actual: %q"
	t.errorf(1, format, a, b)
}

func (t *T) AssertEqual(a, b any) {
	format := "\nexpected: %q\n  actual: %q"
	t.errorf(1, format, a, b)
}

func (t *T) AssertEqualStrings(a, b string) {
	t.Helper()

	if a == b {
		return
	}

	format := "\nexpected: %q\n  actual: %q"
	t.errorf(1, format, a, b)
}

func (t *T) AssertNoError(err error) {
	t.Helper()

	if err != nil {
		t.fatalf(1, "expected no error but got %q", err)
	}
}

func (t *T) AssertEOF(err error) {
	t.Helper()

	if !errors.IsEOF(err) {
		t.fatalf(1, "expected EOF but got %q", err)
	}
}

func (t *T) AssertErrorEquals(expected, actual error) {
	t.Helper()

	if actual == nil {
		t.fatalf(1, "expected %q error but got none", expected)
	}

	if !errors.Is(actual, expected) {
		t.fatalf(1, "expected %q error but got %q", expected, actual)
	}
}

func (t *T) AssertError(err error) {
	t.Helper()

	if err == nil {
		t.fatalf(1, "expected an error but got none")
	}
}
