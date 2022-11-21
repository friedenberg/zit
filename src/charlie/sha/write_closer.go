package sha

import "io"

type WriteCloser interface {
	io.ReaderFrom
	io.WriteCloser
	Sha() Sha
}
