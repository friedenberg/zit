package age

import "io"

type writeCloser struct {
	io.Writer
}

func (w writeCloser) Close() (err error) {
	return
}
