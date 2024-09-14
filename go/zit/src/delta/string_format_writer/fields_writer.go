package string_format_writer

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type Fields struct {
	Boxed   []Field
	Box     bool
	Unboxed []Field
}

func (f Fields) IsBox() bool {
	return f.Box || (len(f.Boxed) > 0 && len(f.Unboxed) > 0)
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
	fields Fields,
) (n int64, err error) {
	if fields.IsBox() {
		if n, err = f.writeStringFormatYesBox(w, fields); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if n, err = f.writeStringFormatNoBox(w, fields); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f *fieldsWriter) writeStringFormatNoBox(
	w interfaces.WriterAndStringWriter,
	fields Fields,
) (n int64, err error) {
	var n1 int64
	var n2 int

	for i, field := range fields.Boxed {
		if i > 0 {
			n2, err = fmt.Fprint(w, " ")
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
	fields Fields,
) (n int64, err error) {
	var n1 int64
	var n2 int

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

	for i, field := range fields.Boxed {
		if i > 0 {
			n2, err = fmt.Fprint(w, " ")
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

	for _, field := range fields.Unboxed {
		n2, err = fmt.Fprint(w, " ")
		n += int64(n2)

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
