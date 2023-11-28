package string_format_writer

import (
	"strings"
	"unicode/utf8"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func MakeRightAligned() schnittstellen.StringFormatWriter[string] {
	return &rightAligned{}
}

type rightAligned struct{}

func (f rightAligned) WriteStringFormat(
	w schnittstellen.WriterAndStringWriter,
	v string,
) (n int64, err error) {
	diff := LenStringMax + 1 - utf8.RuneCountInString(v)

	if diff > 0 {
		v = strings.Repeat(" ", diff-1) + v
	}

	var n1 int

	n1, err = w.WriteString(v)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = w.WriteString(" ")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
