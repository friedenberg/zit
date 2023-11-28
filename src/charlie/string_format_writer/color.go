package string_format_writer

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type color[T any] struct {
	options            ColorOptions
	color              ColorType
	stringFormatWriter schnittstellen.StringFormatWriter[T]
}

func MakeColor[T any](
	o ColorOptions,
	fsw schnittstellen.StringFormatWriter[T],
	c ColorType,
) schnittstellen.StringFormatWriter[T] {
	if o.OffEntirely {
		return fsw
	} else {
		return &color[T]{
			color:              c,
			stringFormatWriter: fsw,
		}
	}
}

func (f *color[T]) WriteStringFormat(
	sw schnittstellen.WriterAndStringWriter,
	e T,
) (n int64, err error) {
	if f.options.OffEntirely {
		return f.stringFormatWriter.WriteStringFormat(sw, e)
	}

	var n1 int

	n1, err = sw.WriteString(string(f.color))
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	var n2 int64
	n2, err = f.stringFormatWriter.WriteStringFormat(sw, e)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = sw.WriteString(string(colorReset))
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
