package test_logz

import (
	"fmt"
	"os"
	"testing"

	"github.com/friedenberg/zit/src/alfa/logz"
)

var (
	Print      = logz.Print
	Printf     = logz.Printf
	PrintDebug = logz.PrintDebug
	LogPrefix  = logz.LogPrefix
)

func Errorf(t *testing.T, format string, args ...interface{}) {
	args = append([]interface{}{LogPrefix()}, args...)
	os.Stderr.WriteString(fmt.Sprintf("%s"+format+"\n", args...))
	t.Fail()
}

func Fatalf(t *testing.T, format string, args ...interface{}) {
	args = append([]interface{}{LogPrefix()}, args...)
	os.Stderr.WriteString(fmt.Sprintf("%s"+format+"\n", args...))
	t.FailNow()
}
