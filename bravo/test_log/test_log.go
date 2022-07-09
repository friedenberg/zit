package test_log

import (
	"fmt"
	"os"
	"testing"

	"github.com/friedenberg/zit/alfa/log"
)

var (
	Print      = log.Print
	Printf     = log.Printf
	PrintDebug = log.PrintDebug
	LogPrefix  = log.LogPrefix
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
