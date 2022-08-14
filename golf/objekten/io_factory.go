package objekten

import (
	"io"

	"github.com/friedenberg/zit/delta/age_io"
)

func (s Store) ReadCloser(p string) (r io.ReadCloser, err error) {
	return objekte.NewFileReader(s.Age, p)
}

func (s Store) WriteCloser(p string) (w io.WriteCloser, err error) {
	return objekte.NewWriterMoverPrenamed(s.Age, p)
}
