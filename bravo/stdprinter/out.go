package stdprinter

import (
	"fmt"
	"os"
)

func Outf(f string, a ...interface{}) {
	printerChannel <- printerLine{
		file: os.Stdout,
		line: fmt.Sprintf(f, a...),
	}
}

func Out(a ...interface{}) {
	printerChannel <- printerLine{
		file: os.Stdout,
		line: fmt.Sprintln(a...),
	}
}
