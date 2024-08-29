package string_format_writer

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type Field struct {
	Key, Value         string
	DisableValueQuotes bool
	ColorType
	Prefix    string
	Separator rune
}

type formatCliField struct {
	ColorOptions
	truncate CliFormatTruncation
}

func MakeCliFormatField(
	truncate CliFormatTruncation,
	co ColorOptions,
) *formatCliField {
	return &formatCliField{
		truncate:     truncate,
		ColorOptions: co,
	}
}

func (f *formatCliField) WriteStringFormat(
	w interfaces.WriterAndStringWriter,
	field Field,
) (n int64, err error) {
	var n1 int

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
