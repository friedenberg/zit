package alfred

import (
	"bufio"
	"encoding/json"
	"io"

	"github.com/friedenberg/zit/bravo/errors"
)

type Writer interface {
	io.WriteCloser
	WriteItem(i Item) (n int, err error)
}

type writer struct {
	wBuf            *bufio.Writer
	afterFirstWrite bool
}

func NewWriter(out io.Writer) (w *writer, err error) {
	w = &writer{
		wBuf: bufio.NewWriter(out),
	}

	if err = w.open(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (w writer) open() (err error) {
	_, err = w.wBuf.WriteString("{\"items\":[\n")

	return
}

func (w *writer) setAfterFirstWrite() {
	w.afterFirstWrite = true
}

func (w *writer) WriteItem(i Item) (n int, err error) {
	var b []byte

	if b, err = json.Marshal(i); err != nil {
		err = errors.Error(err)
		return
	}

	if n, err = w.Write(b); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (w *writer) Write(p []byte) (n int, err error) {
	defer w.setAfterFirstWrite()

	var n1 int

	if w.afterFirstWrite {
		n1, err = w.wBuf.WriteString(",")

		n += n1
	}

	n1, err = w.wBuf.Write(p)

	n += 1

	return
}

func (w writer) Close() (err error) {
	if _, err = w.wBuf.WriteString("]}\n"); err != nil {
		err = errors.Error(err)
		return
	}

	if err = w.wBuf.Flush(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
