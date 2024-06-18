package test_logz

import (
	"fmt"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

type TC struct {
	T
	StackInfo
}

func (t *TC) ui(args ...interface{}) {
	errors.SetTesting()
	args = append([]interface{}{t.StackInfo}, args...)
	fmt.Fprintln(os.Stderr, args...)
}

func (t *TC) logf(format string, args ...interface{}) {
	errors.SetTesting()
	args = append([]interface{}{t.StackInfo}, args...)
	fmt.Fprintf(os.Stderr, "%s"+format+"\n", args...)
}

func (t *TC) errorf(format string, args ...interface{}) {
	t.logf(format, args...)
	t.Fail()
}

func (t *TC) fatalf(format string, args ...interface{}) {
	t.logf(format, args...)
	t.FailNow()
}

func (t *TC) Log(args ...interface{}) {
	t.ui(args...)
}

func (t *TC) Logf(format string, args ...interface{}) {
	t.logf(format, args...)
}

func (t *TC) Errorf(format string, args ...interface{}) {
	t.Helper()
	t.errorf(format, args...)
}

func (t *TC) Fatalf(format string, args ...interface{}) {
	t.Helper()
	t.fatalf(format, args...)
}

// TODO-P3 move to AssertNotEqual
func (t *TC) NotEqual(a, b any) {
	format := "\nexpected: %q\n  actual: %q"
	t.errorf(format, a, b)
}

func (t *TC) AssertEqual(a, b any) {
	format := "\nexpected: %q\n  actual: %q"
	t.errorf(format, a, b)
}

func (t *TC) AssertEqualStrings(a, b string) {
	t.Helper()

	if a == b {
		return
	}

	format := "\nexpected: %q\n  actual: %q"
	t.errorf(format, a, b)
}

func (t *TC) AssertNoError(err error) {
	t.Helper()

	if err != nil {
		t.fatalf("expected no error but got %q", err)
	}
}

func (t *TC) AssertEOF(err error) {
	t.Helper()

	if !errors.IsEOF(err) {
		t.fatalf("expected EOF but got %q", err)
	}
}

func (t *TC) AssertErrorEquals(expected, actual error) {
	t.Helper()

	if actual == nil {
		t.fatalf("expected %q error but got none", expected)
	}

	if !errors.Is(actual, expected) {
		t.fatalf("expected %q error but got %q", expected, actual)
	}
}

func (t *TC) AssertError(err error) {
	t.Helper()

	if err == nil {
		t.fatalf("expected an error but got none")
	}
}
