package string_format_writer

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

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
	fields []Field,
) (n int64, err error) {
	var n1 int64
	var n2 int

	for i, field := range fields {
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
