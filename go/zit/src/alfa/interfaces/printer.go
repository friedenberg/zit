package interfaces

import (
	"io"
	"os"
)

type Printer interface {
	io.Writer
	GetPrinter() Printer

	GetFile() *os.File
	IsTty() bool
	Print(a ...interface{}) (err error)
	Printf(f string, a ...interface{}) (err error)
}
