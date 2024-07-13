package test_object_metadata_io

import (
	"io"
)

type nopReadWriteCloser struct {
	io.ReadWriter
}

func NopReadWriteCloser(rw io.ReadWriter) *nopReadWriteCloser {
	return &nopReadWriteCloser{
		ReadWriter: rw,
	}
}

func (nrwc *nopReadWriteCloser) Close() (err error) {
	return
}

func (b nopReadWriteCloser) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = io.Copy(b.ReadWriter, r)
	return
}

func (b nopReadWriteCloser) WriteTo(w io.Writer) (n int64, err error) {
	n, err = io.Copy(w, b.ReadWriter)
	return
}
