package verzeichnisse

import (
	"io"

	"github.com/friedenberg/zit/bravo/sha"
)

type IdTransformer func(sha.Sha) string

type Reader interface {
	Begin() (err error)
	ReadRow(string, Row) (err error)
	End() (err error)
}

type ReadCloserFactory interface {
	ReadCloser(string) (io.ReadCloser, error)
}

type WriteCloserFactory interface {
	WriteCloser(string) (io.WriteCloser, error)
}
