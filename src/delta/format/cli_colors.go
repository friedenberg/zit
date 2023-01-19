package format

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type color string
type ColorType color

const (
	colorReset          = "\u001b[0m"
	colorBlack          = "\u001b[30m"
	colorRed            = "\u001b[31m"
	colorGreen          = "\u001b[32m"
	colorYellow         = "\u001b[33m"
	colorBlue           = "\u001b[34m"
	colorMagenta        = "\u001b[35m"
	colorCyan           = "\u001b[36m"
	colorWhite          = "\u001b[37m"
	colorItalic         = "\u001b[3m"
	colorNone           = ""
	ColorTypePointer    = colorBlue
	ColorTypeConstant   = colorItalic
	ColorTypeType       = colorYellow
	ColorTypeIdentifier = colorCyan
	ColorTypeTitle      = colorRed
)

func MakeFormatWriterNoopColor(
	wf schnittstellen.FuncWriter,
	c ColorType,
) schnittstellen.FuncWriter {
	return wf
}

func MakeFormatWriterWithColor(
	wf schnittstellen.FuncWriter,
	c ColorType,
) schnittstellen.FuncWriter {
	return func(w io.Writer) (n int64, err error) {
		return Write(
			w,
			MakeFormatString(string(c)),
			wf,
			MakeFormatString(string(colorReset)),
		)
	}
}
