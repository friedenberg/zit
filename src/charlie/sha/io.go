package sha

import "io"

type ReadCloser interface {
	io.WriterTo
	io.ReadCloser
	Sha() Sha
}

type WriteCloser interface {
	io.ReaderFrom
	io.WriteCloser
	Sha() Sha
}
