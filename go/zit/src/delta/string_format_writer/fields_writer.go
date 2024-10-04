package string_format_writer

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type Box struct {
	Header              []Field
	Contents            []Field
	Trailer             []Field
	EachFieldOnANewline bool
}

func (b Box) GetSeparator() string {
	if b.EachFieldOnANewline && false {
		return "\n"
	} else {
		return " "
	}
}

type fieldsWriter struct {
	fieldWriter
}

func MakeCliFormatFields(
	truncate CliFormatTruncation,
	co ColorOptions,
) *fieldsWriter {
	return &fieldsWriter{
		fieldWriter: fieldWriter{
			truncate:     truncate,
			ColorOptions: co,
		},
	}
}

func (f *fieldsWriter) WriteStringFormat(
	w interfaces.WriterAndStringWriter,
	fields Box,
) (n int64, err error) {
	if n, err = f.writeStringFormatYesBox(w, fields); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *fieldsWriter) writeStringFormatNoBox(
	w interfaces.WriterAndStringWriter,
	fields Box,
) (n int64, err error) {
	var n1 int64
	var n2 int

	separator := fields.GetSeparator()

	for i, field := range fields.Contents {
		if i > 0 {
			n2, err = fmt.Fprint(w, separator)
			n += int64(n2)

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		n1, err = f.fieldWriter.WriteStringFormat(w, field)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f *fieldsWriter) writeStringFormatYesBox(
	w interfaces.WriterAndStringWriter,
	fields Box,
) (n int64, err error) {
	var n1 int64
	var n2 int

	separator := fields.GetSeparator()

	for _, field := range fields.Header {
		n1, err = f.fieldWriter.WriteStringFormat(w, field)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = fmt.Fprint(w, separator)
		n += int64(n2)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	n1, err = f.fieldWriter.WriteStringFormat(
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
		if i > 0 && !field.NeedsNewline {
			n2, err = fmt.Fprint(w, separator)
			n += int64(n2)

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		n1, err = f.fieldWriter.WriteStringFormat(w, field)
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

	n1, err = f.fieldWriter.WriteStringFormat(
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

		n1, err = f.fieldWriter.WriteStringFormat(w, field)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
