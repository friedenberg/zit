package string_format_writer

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type Field struct {
	ColorType
	Separator          rune
	Key, Value         string
	DisableValueQuotes bool
	NoTruncate         bool
	NeedsNewline       bool
}

type BoxHeader struct {
	Value        string
	RightAligned bool
}

type HeaderWriter[T any] interface {
	WriteBoxHeader(*BoxHeader, T) error
}

type Box struct {
	Header                   BoxHeader
	Contents                 []Field
	Trailer                  []Field
	EachFieldOnANewline      bool
	IdFieldsSeparatedByLines bool
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

func (f *fieldsWriter) EncodeStringTo(
	box Box,
	w interfaces.WriterAndStringWriter,
) (n int64, err error) {
	var n1 int64
	var n2 int

	separatorSameLine := " "
	separatorNextLine := "\n" + StringIndentWithSpace

	if box.Header.Value != "" {
		headerWriter := w

		if box.Header.RightAligned {
			headerWriter = rightAligned2{w}
		}

		n2, err = headerWriter.WriteString(box.Header.Value)
		n += int64(n2)

		if err != nil {
			err = errors.Wrapf(err, "Headers: %#v", box.Header)
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

	for i, field := range box.Contents {
		if i > 0 {
			if field.NeedsNewline {
				n2, err = w.WriteString(separatorNextLine)
				n += int64(n2)

				if err != nil {
					err = errors.Wrap(err)
					return
				}
			} else {
				n2, err = fmt.Fprint(w, separatorSameLine)
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

	if separatorSameLine == "\n" {
		n2, err = w.WriteString(separatorSameLine)
		n += int64(n2)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	closingBracket := "]"

	if len(box.Trailer) > 0 && false {
		closingBracket = "\n" + StringIndent + " ]"
	}

	n1, err = f.writeStringFormatField(
		w,
		Field{
			Value:     closingBracket,
			ColorType: ColorTypeNormal,
		},
	)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, field := range box.Trailer {
		n2, err = fmt.Fprint(w, separatorSameLine)
		n += int64(n2)

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

	if !field.NoTruncate && trunc > 0 && len(field.Value) > int(trunc) {
		field.Value = field.Value[:trunc+1]
		ellipsis = "…"
	}

	format := "%s%s%s%s"

	if (strings.ContainsRune(field.Value, ' ') || field.ColorType == ColorTypeUserData) &&
		!field.DisableValueQuotes {
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
