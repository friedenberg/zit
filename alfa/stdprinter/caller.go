package stdprinter

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

type callerFormatter struct {
	*strings.Builder
}

func (t callerFormatter) Print(args ...interface{}) {
	t.WriteString(fmt.Sprintln(args...))
}

func (t callerFormatter) Printf(format string, args ...interface{}) {
	t.WriteString(fmt.Sprintf(format, args...))
}

func (t callerFormatter) Detail() bool {
	return true
}

func Caller(i int, msg ...interface{}) {
	cf := callerFormatter{&strings.Builder{}}

	cf.Print(msg...)
	xerrors.Caller(i + 1).Format(cf)

	Err(cf.String())
}
