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
	Skip int
}

func (t T) NotEqual(a, b any) {
	format := "\nexpected: %v\n  actual: %v"
	args := []any{a, b}

	errors.SetTesting()
	si, _ := MakeStackInfo(t.Skip + 1)
	args = append([]interface{}{si}, args...)
	os.Stderr.WriteString(fmt.Sprintf("%s"+format+"\n", args...))
	t.Fail()
}

func (t T) Errorf(format string, args ...interface{}) {
	errors.SetTesting()
	si, _ := MakeStackInfo(t.Skip + 1)
	args = append([]interface{}{si}, args...)
	os.Stderr.WriteString(fmt.Sprintf("%s"+format+"\n", args...))
	t.Fail()
}

func (t T) Fatalf(format string, args ...interface{}) {
	errors.SetTesting()
	si, _ := MakeStackInfo(t.Skip + 1)
	args = append([]interface{}{si}, args...)
	os.Stderr.WriteString(fmt.Sprintf("%s"+format+"\n", args...))
	t.FailNow()
}
