package sha

import "io"

type WriteCloser interface {
	io.WriteCloser
	Sha() Sha
}
