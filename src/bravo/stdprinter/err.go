package stdprinter

import (
	"fmt"
	"os"
)

func PanicIfError(err interface{}) {
	if err == nil {
		return
	}

	switch t := err.(type) {
	case func() error:
		PanicIfError(t())
	case error:
		panic(err)
	}
}

func Error(err error) {
	if err == nil {
		return
	}

	printerChannel <- printerLine{
		file: os.Stderr,
		line: fmt.Sprintf("%+v", err),
	}
}

func Errf(f string, a ...interface{}) {
	printerChannel <- printerLine{
		file: os.Stderr,
		line: fmt.Sprintf(f, a...),
	}
}

func Err(a ...interface{}) {
	printerChannel <- printerLine{
		file: os.Stderr,
		line: fmt.Sprintln(a...),
	}
}
