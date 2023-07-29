package format

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type (
	color     string
	ColorType color
)

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

type ColorOptions struct {
	OffEntirely bool
}

type colorStringFormat[T any] struct {
	options            ColorOptions
	color              ColorType
	stringFormatWriter schnittstellen.StringFormatWriter[T]
}

func MakeColorStringFormatWriter[T any](
	o ColorOptions,
	fsw schnittstellen.StringFormatWriter[T],
	c ColorType,
) schnittstellen.StringFormatWriter[T] {
	if o.OffEntirely {
		return fsw
	} else {
		return &colorStringFormat[T]{
			color:              c,
			stringFormatWriter: fsw,
		}
	}
}

func (f *colorStringFormat[T]) WriteStringFormat(
	sw io.StringWriter,
	e T,
) (n int64, err error) {
	if f.options.OffEntirely {
		return f.stringFormatWriter.WriteStringFormat(sw, e)
	}

	var n1 int

	n1, err = sw.WriteString(string(f.color))
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	var n2 int64
	n2, err = f.stringFormatWriter.WriteStringFormat(sw, e)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = sw.WriteString(string(colorReset))
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeStringerStringFormatWriter[T schnittstellen.Stringer]() schnittstellen.StringFormatWriter[T] {
	return &stringerFormatWriter[T]{}
}

type stringerFormatWriter[T schnittstellen.Stringer] struct{}

func (f *stringerFormatWriter[T]) WriteStringFormat(
	sw io.StringWriter,
	e T,
) (n int64, err error) {
	var n1 int

	n1, err = sw.WriteString(e.String())
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type stringWriterFormatString[T ~string] struct{}

func MakeStringStringFormatWriter[T ~string]() schnittstellen.StringFormatWriter[T] {
	return &stringWriterFormatString[T]{}
}

func (f *stringWriterFormatString[T]) WriteStringFormat(
	sw io.StringWriter,
	e T,
) (n int64, err error) {
	var n1 int

	n1, err = sw.WriteString(string(e))
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
