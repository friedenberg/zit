package objekte

import "io"

type ReadCloserFactory interface {
	ReadCloser(string) (io.ReadCloser, error)
}

type WriteCloserFactory interface {
	WriteCloser(string) (io.WriteCloser, error)
}
