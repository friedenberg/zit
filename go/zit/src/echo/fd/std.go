package fd

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
)

type Std struct {
	ui.Printer
}

func MakeStd(f *os.File) Std {
	return Std{
		Printer: ui.MakePrinter(f),
	}
}
