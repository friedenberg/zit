package format

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

const (
	StringDRArrow        = "↳"
	StringNew            = "new"
	StringSame           = "same"
	StringChanged        = "changed"
	StringDeleted        = "deleted"
	StringUpdated        = "updated"
	StringArchived       = "archived"
	StringUnchanged      = "unchanged"
	StringRecognized     = "recognized"
	StringCheckedOut     = "checked out"
	StringWouldDelete    = "would delete"
	StringUnrecognized   = "unrecognized"
	StringFormatDateTime = "06-01-02 15:04:05"
	StringIndent         = "                 "
	LenStringMax         = len(StringIndent) // TODO-P4 use reflection?
)

func MakeBracketWrappedStringFormatWriter[T any](
	sfw schnittstellen.StringFormatWriter[T],
) schnittstellen.StringFormatWriter[T] {
	return &bracketWrappedStringFormatWriter[T]{
		stringFormatWriter: sfw,
	}
}

type bracketWrappedStringFormatWriter[T any] struct {
	stringFormatWriter schnittstellen.StringFormatWriter[T]
}

func (f bracketWrappedStringFormatWriter[T]) WriteStringFormat(
	w io.StringWriter,
	e T,
) (n int64, err error) {
	var (
		n1 int
		n2 int64
	)

	n1, err = w.WriteString("[")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = f.stringFormatWriter.WriteStringFormat(w, e)
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = w.WriteString("]")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeRightAlignedStringFormatWriter() schnittstellen.StringFormatWriter[string] {
	return &stringFormatWriterRightAligned{}
}

type stringFormatWriterRightAligned struct{}

func (f stringFormatWriterRightAligned) WriteStringFormat(
	w io.StringWriter,
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

func MakeFormatStringRightAligned(
	f string,
	args ...any,
) schnittstellen.FuncWriter {
	return func(w io.Writer) (n int64, err error) {
		f = fmt.Sprintf(f+" ", args...)

		diff := LenStringMax + 1 - utf8.RuneCountInString(f)

		if diff > 0 {
			f = strings.Repeat(" ", diff) + f
		}

		var n1 int

		if n1, err = io.WriteString(w, f); err != nil {
			n = int64(n1)
			err = errors.Wrap(err)
			return
		}

		n = int64(n1)

		return
	}
}

func MakeWriterFormatStringIndentedHeader(
	cw FuncColorWriter,
	indentString string,
) schnittstellen.FuncWriterFormat[string] {
	return func(w io.Writer, v string) (n int64, err error) {
		return Write(
			w,
			MakeFormatString(indentString),
			cw(MakeFormatString("%s", v), ColorTypeTitle),
		)
	}
}
