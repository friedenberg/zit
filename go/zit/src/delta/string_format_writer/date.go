package string_format_writer

import (
	"io"
	"time"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type Clock interface {
	GetTime() time.Time
}

type Date[T any] struct {
	Clock
	Format string
	interfaces.StringFormatWriter[T]
}

func MakeDefaultDatePrefixFormatWriter[T any](
	clock Clock,
	f interfaces.StringFormatWriter[T],
) interfaces.StringFormatWriter[T] {
	return &Date[T]{
		Clock:              clock,
		Format:             StringFormatDateTime,
		StringFormatWriter: f,
	}
}

func (f *Date[T]) WriteStringFormat(
	e T,
	w interfaces.WriterAndStringWriter,
) (n int64, err error) {
	d := f.GetTime().Format(f.Format)

	var n1 int

	n1, err = io.WriteString(w, d)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = io.WriteString(w, " ")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	var n2 int64

	n2, err = f.StringFormatWriter.WriteStringFormat(e, w)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
