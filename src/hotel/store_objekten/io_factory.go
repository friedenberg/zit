package objekten

import (
	"io"

	age_io "github.com/friedenberg/zit/src/echo/age_io"
)

func (s Store) ReadCloser(p string) (r io.ReadCloser, err error) {
	o := age_io.FileReadOptions{
		Age:  s.Age,
		Path: p,
	}

	return age_io.NewFileReader(o)
}

func (s Store) WriteCloser(p string) (w io.WriteCloser, err error) {
	return age_io.NewMover(
		age_io.MoveOptions{
			Age:       s.Age,
			FinalPath: p,
		},
	)
}
