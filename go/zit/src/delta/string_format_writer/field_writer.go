package string_format_writer

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type Field struct {
	Key, Value string
	ColorType
	Prefix    string
	Separator rune
}

type formatCliField struct {
	ColorOptions
	truncate int
}

func MakeCliFormatField(
	truncate int,
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

	format := "%s%s%s%s"

	if field.ColorType == ColorTypeUserData {
		format = "%s%q%s%s"

		if f.truncate > 0 && len(field.Value) > f.truncate {
			field.Value = field.Value[:f.truncate+1]
			ellipsis = "…"
		}
	}

	n1, err = fmt.Fprintf(w, format, preColor, field.Value, postColor, ellipsis)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
