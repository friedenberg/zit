package objekten

import (
	"io"

	age_io "github.com/friedenberg/zit/delta/age_io"
)

func (s Store) ReadCloser(p string) (r io.ReadCloser, err error) {
	return age_io.NewFileReader(s.Age, p)
}

func (s Store) WriteCloser(p string) (w io.WriteCloser, err error) {
	return age_io.NewWriterMoverPrenamed(s.Age, p)
}
