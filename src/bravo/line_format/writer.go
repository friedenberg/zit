package line_format

import (
	"bufio"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Writer []string

func NewWriter() *Writer {
	w := Writer(make([]string, 0))
	return &w
}

func (w *Writer) WriteTo(out io.Writer) (n int64, err error) {
	w1 := bufio.NewWriter(out)
	defer func() {
		if err == nil {
			err = w1.Flush()
		}
	}()

	var n1 int

	for _, l := range *w {
		n1, err = w1.WriteString(fmt.Sprintln(l))
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (w *Writer) WriteExactlyOneEmpty() {
	if len(*w) == 0 {
		w.WriteEmpty()
		return
	}

	if (*w)[len(*w)-1] != "" {
		w.WriteEmpty()
		return
	}

	return
}

func (w *Writer) WriteEmpty() {
	*w = append(*w, "")
}

func (w *Writer) WriteLines(ls ...string) {
	*w = append(*w, ls...)
}

func (w *Writer) WriteStringers(ss ...fmt.Stringer) {
	for _, s := range ss {
		w.WriteLines(s.String())
	}
}

func (w *Writer) WriteFormat(f string, values ...interface{}) {
	w.WriteLines(fmt.Sprintf(f, values...))
}

func (w *Writer) WriteFormats(f string, values ...interface{}) {
	for _, v := range values {
		w.WriteLines(fmt.Sprintf(f, v))
	}
}
