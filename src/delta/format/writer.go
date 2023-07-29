package format

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

// func StringSep(
// 	delim byte,
// 	sfs ...fmt.Stringer,
// ) string {
//   sb := &strings.Builder{}

// 	for _, sf := range sfs {
// 		var n1 int

// 		if n1, err = io.WriteString(w, sf.String()); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}

// 		n += int64(n1)
// 	}

// 	return sb.String()
// }

// func WriteStringSep(
// 	delim byte,
// 	w1 io.Writer,
// 	sfs ...fmt.Stringer,
// ) (n int64, err error) {
// 	w := bufio.NewWriter(w1)
// 	defer errors.DeferredFlusher(&err, w)

// 	for _, sf := range sfs {
// 		var n1 int

// 		if n1, err = io.WriteString(w, sf.String()); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}

// 		n += int64(n1)
// 	}

// 	return
// }

func Write(
	w io.Writer,
	wffs ...schnittstellen.FuncWriter,
) (n int64, err error) {
	for _, wf := range wffs {
		var n1 int64

		if n1, err = wf(w); err != nil {
			err = errors.Wrap(err)
			return
		}

		n += n1
	}

	return
}

func WriteLines(
	w1 io.Writer,
	wffs ...schnittstellen.FuncWriter,
) (n int64, err error) {
	w := bufio.NewWriter(w1)
	defer errors.DeferredFlusher(&err, w)

	sb := &strings.Builder{}

	for _, wf := range wffs {
		sb.Reset()

		if _, err = wf(sb); err != nil {
			err = errors.Wrap(err)
			return
		}

		sb.WriteByte('\n')

		var n1 int

		if n1, err = io.WriteString(w, sb.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		n += int64(n1)
	}

	return
}

// TODO-P2 rename
func MakeWriterTo[T any](
	w io.Writer,
	wf schnittstellen.FuncWriterFormat[T],
) func(T) error {
	return func(e T) (err error) {
		if _, err = wf(w, e); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func MakeWriterToWithNewLines[T any](
	w1 io.Writer,
	wf schnittstellen.FuncWriterFormat[T],
) func(T) error {
	w := bufio.NewWriter(w1)

	return collections.MakeSyncSerializer(
		func(e T) (err error) {
			errors.TodoP3("modify flushing behavior based on w1 being a TTY")
			defer errors.DeferredFlusher(&err, w)

			if _, err = wf(w, e); err != nil {
				err = errors.Wrap(err)
				return
			}

			if _, err = io.WriteString(w, "\n"); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)
}

func MakeWriterToWithNewLinesPtr[T any](
	w1 io.Writer,
	wf schnittstellen.FuncWriterFormat[T],
) func(*T) error {
	w := bufio.NewWriter(w1)

	return collections.MakeSyncSerializer(
		func(e *T) (err error) {
			errors.TodoP3("modify flushing behavior based on w1 being a TTY")
			defer errors.DeferredFlusher(&err, w)

			if _, err = wf(w, *e); err != nil {
				err = errors.Wrap(err)
				return
			}

			if _, err = io.WriteString(w, "\n"); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)
}

func MakeFuncStringFormatWriter[T any](
	f schnittstellen.FuncStringWriterFormat[T],
) schnittstellen.StringFormatWriter[T] {
	return stringFormatWriterFunc[T](f)
}

type stringFormatWriterFunc[T any] schnittstellen.FuncStringWriterFormat[T]

func (f stringFormatWriterFunc[T]) WriteStringFormat(
	w io.StringWriter,
	e T,
) (int64, error) {
	return schnittstellen.FuncStringWriterFormat[T](f)(w, e)
}

func MakeDelimFuncStringFormatWriter[T any](
	delim string,
	w1 schnittstellen.WriterAndStringWriter,
	f schnittstellen.StringFormatWriter[T],
) func(T) error {
	w := bufio.NewWriter(w1)

	return collections.MakeSyncSerializer(
		func(e T) (err error) {
			errors.TodoP3("modify flushing behavior based on w1 being a TTY")
			defer errors.DeferredFlusher(&err, w)

			if _, err = f.WriteStringFormat(w, e); err != nil {
				err = errors.Wrap(err)
				return
			}

			if _, err = io.WriteString(w, delim); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)
}
