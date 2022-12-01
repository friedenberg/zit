package format

import (
	"bufio"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

func Write(
	w io.Writer,
	wffs ...WriterFunc,
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

func MakeWriterTo[T any](
	w io.Writer,
	wf FormatWriterFunc[T],
) func(*T) error {
	return func(e *T) (err error) {
		if _, err = wf(w, e); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func MakeWriterToWithNewLines[T any](
	w1 io.Writer,
	wf FormatWriterFunc[T],
) func(*T) error {

	w := bufio.NewWriter(w1)

	return func(e *T) (err error) {
		//TODO modify flushing behavior based on w1 being a TTY
		defer errors.Deferred(&err, w.Flush)

		if _, err = wf(w, e); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = io.WriteString(w, "\n"); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

type Writer struct {
	lastWasNewline bool
	elements       []FuncWriter
}

func NewWriter() *Writer {
	w := &Writer{
		elements: make([]FuncWriter, 0),
	}

	return w
}

func (w *Writer) WriteTo(out io.Writer) (n int64, err error) {
	w1 := bufio.NewWriter(out)
	defer errors.Deferred(&err, w1.Flush)

	var n1 int64
	var n2 int

	for _, l := range w.elements {
		if n1, err = l(w1); err != nil {
			err = errors.Wrap(err)
			return
		}

		n += n1

		if n2, err = io.WriteString(w1, "\n"); err != nil {
			err = errors.Wrap(err)
			return
		}

		n += int64(n2)
	}

	return
}

func (w *Writer) WriteExactlyOneEmpty() {
	if len(w.elements) == 0 || !w.lastWasNewline {
		w.WriteEmpty()
		return
	}

	return
}

func (w *Writer) WriteEmpty() {
	w.lastWasNewline = true

	w.elements = append(
		w.elements,
		func(_ io.Writer) (_ int64, _ error) {
			return
		},
	)
}

func (w *Writer) WriteLines(ls ...string) {
	w.lastWasNewline = false

	for _, v := range ls {
		w.elements = append(
			w.elements,
			MakeFormatString("%s", v),
		)
	}
}

func (w *Writer) WriteStringers(ss ...fmt.Stringer) {
	w.lastWasNewline = false

	for _, v := range ss {
		w.elements = append(
			w.elements,
			MakeStringer(v),
		)
	}
}

func (w *Writer) WriteFormat(f string, values ...interface{}) {
	w.lastWasNewline = false

	w.elements = append(
		w.elements,
		MakeFormatString(f, values...),
	)
}

func (w *Writer) WriteFormats(f string, values ...interface{}) {
	w.lastWasNewline = false

	for _, v := range values {
		w.elements = append(
			w.elements,
			MakeFormatString(f, v),
		)
	}
}
