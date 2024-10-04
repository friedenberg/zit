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

type Box struct {
	Header              []Field
	Contents            []Field
	Trailer             []Field
	EachFieldOnANewline bool
}

func (b Box) GetSeparator() string {
	return " "
	if b.EachFieldOnANewline || true {
		return "\n"
	} else {
		return " "
	}
}

type fieldsWriter struct {
	ColorOptions
	truncate CliFormatTruncation
	rightAligned
}

func MakeCliFormatFields(
	truncate CliFormatTruncation,
	co ColorOptions,
) *fieldsWriter {
	return &fieldsWriter{
		truncate:     truncate,
		ColorOptions: co,
	}
}

func (f *fieldsWriter) WriteStringFormat(
	w interfaces.WriterAndStringWriter,
	fields Box,
) (n int64, err error) {
	var n1 int64
	var n2 int

	separator := fields.GetSeparator()

	for _, field := range fields.Header {
		n1, err = f.writeStringFormatField(rightAligned2{w}, field)
		n += n1

		if err != nil {
			err = errors.Wrapf(err, "Headers: %#v", fields.Header)
			return
		}
	}

	n1, err = f.writeStringFormatField(
		w,
		Field{
			Value:     "[",
			ColorType: ColorTypeNormal,
		},
	)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	for i, field := range fields.Contents {
		if i > 0 {
			if field.NeedsNewline {
				n2, err = w.WriteString("\n")
				n += int64(n2)

				if err != nil {
					err = errors.Wrap(err)
					return
				}
			} else {
				n2, err = fmt.Fprint(w, separator)
				n += int64(n2)

				if err != nil {
					err = errors.Wrap(err)
					return
				}
			}
		}

		n1, err = f.writeStringFormatField(w, field)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if separator == "\n" {
		n2, err = w.WriteString(separator)
		n += int64(n2)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	n1, err = f.writeStringFormatField(
		w,
		Field{
			Value:     "]",
			ColorType: ColorTypeNormal,
		},
	)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	for i, field := range fields.Trailer {
		if i == 0 {
			// n2, err = fmt.Fprint(w, " ")
		} else {
			n2, err = fmt.Fprint(w, separator)
			n += int64(n2)
		}

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n1, err = f.writeStringFormatField(w, field)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f *fieldsWriter) writeStringFormatField(
	w interfaces.WriterAndStringWriter,
	field Field,
) (n int64, err error) {
	var n1 int

	if field.Prefix != "" {
		n1, err = w.WriteString(field.Prefix)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
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
