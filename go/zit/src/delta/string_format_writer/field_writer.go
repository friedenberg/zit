package string_format_writer

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type Field struct {
	NeedsNewline bool
	Prefix       string
	ColorType
	Separator          rune
	Key, Value         string
	DisableValueQuotes bool
}

type fieldWriter struct {
	ColorOptions
	truncate CliFormatTruncation
}

func MakeCliFormatField(
	truncate CliFormatTruncation,
	co ColorOptions,
) *fieldWriter {
	return &fieldWriter{
		truncate:     truncate,
		ColorOptions: co,
	}
}

func (f *fieldWriter) WriteStringFormat(
	w interfaces.WriterAndStringWriter,
	field Field,
) (n int64, err error) {
	// ui.Debug().Printf("%#v", field)
	var n1 int

	if field.NeedsNewline {
		n1, err = w.WriteString("\n")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	n1, err = w.WriteString(field.Prefix)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if field.Key != "" {
		if field.Separator == '\x00' {
			field.Separator = '='
		}

		n1, err = fmt.Fprintf(w, "%s%c", field.Key, field.Separator)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	preColor, postColor, ellipsis := field.ColorType, colorReset, ""

	if f.OffEntirely {
		preColor, postColor = "", ""
	}

	trunc := f.truncate

	if trunc == CliFormatTruncation66CharEllipsis {
		trunc = 66
	}

	if trunc > 0 && len(field.Value) > int(trunc) {
		field.Value = field.Value[:trunc+1]
		ellipsis = "…"
	}

	format := "%s%s%s%s"

	if field.ColorType == ColorTypeUserData && !field.DisableValueQuotes {
		format = "%s%q%s%s"
	}

	n1, err = fmt.Fprintf(w, format, preColor, field.Value, postColor, ellipsis)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
