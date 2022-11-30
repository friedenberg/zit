package sha

import "io"

type nopReadCloser struct {
	io.ReadCloser
}

func MakeNopReadCloser(rc io.ReadCloser) ReadCloser {
	return nopReadCloser{
		ReadCloser: rc,
	}
}

func (nrc nopReadCloser) WriteTo(w io.Writer) (n int64, err error) {
	return io.Copy(w, nrc.ReadCloser)
}

func (nrc nopReadCloser) Sha() Sha {
	return Sha{}
}
