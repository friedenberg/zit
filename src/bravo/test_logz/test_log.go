package test_logz

import (
	"fmt"
	"os"
	"testing"

	"github.com/friedenberg/zit/src/alfa/errors"
)

var (
	Print         = errors.Print
	Printf        = errors.Printf
	PrintDebug    = errors.PrintDebug
	MakeStackInfo = errors.MakeStackInfo
)

type (
	StackInfo = errors.StackInfo
)

type T struct {
	*testing.T
	Skip int
}

func (t T) Errorf(format string, args ...interface{}) {
	si, _ := MakeStackInfo(t.Skip + 1)
	args = append([]interface{}{si}, args...)
	os.Stderr.WriteString(fmt.Sprintf("%s"+format+"\n", args...))
	t.Fail()
}

func (t T) Fatalf(format string, args ...interface{}) {
	si, _ := MakeStackInfo(t.Skip + 1)
	args = append([]interface{}{si}, args...)
	os.Stderr.WriteString(fmt.Sprintf("%s"+format+"\n", args...))
	t.FailNow()
}
