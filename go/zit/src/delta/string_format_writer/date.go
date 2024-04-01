package string_format_writer

import (
	"io"
	"time"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
)

type Clock interface {
	GetTime() time.Time
}

type Date[T any] struct {
	Clock
	Format string
	schnittstellen.StringFormatWriter[T]
}

func MakeDefaultDatePrefixFormatWriter[T any](
	clock Clock,
	f schnittstellen.StringFormatWriter[T],
) schnittstellen.StringFormatWriter[T] {
	return &Date[T]{
		Clock:              clock,
		Format:             StringFormatDateTime,
		StringFormatWriter: f,
	}
}

func (f *Date[T]) WriteStringFormat(
	w schnittstellen.WriterAndStringWriter,
	e T,
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

	n2, err = f.StringFormatWriter.WriteStringFormat(w, e)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
