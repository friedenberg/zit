package store_objekten

import (
	"io"

	"github.com/friedenberg/zit/src/echo/age_io"
)

func (s Store) ReadCloserObjekten(p string) (io.ReadCloser, error) {
	o := age_io.FileReadOptions{
		Age:  s.Age,
		Path: p,
	}

	return age_io.NewFileReader(o)
}

func (s Store) ReadCloserVerzeichnisse(p string) (io.ReadCloser, error) {
	o := age_io.FileReadOptions{
		Age:  s.Age,
		Path: p,
	}

	return age_io.NewFileReader(o)
}

func (s Store) WriteCloserObjekten(p string) (w io.WriteCloser, err error) {
	return age_io.NewMover(
		age_io.MoveOptions{
			Age:       s.Age,
			FinalPath: p,
			LockFile:  true,
		},
	)
}

func (s Store) WriteCloserVerzeichnisse(p string) (w io.WriteCloser, err error) {
	return age_io.NewMover(
		age_io.MoveOptions{
			Age:       s.Age,
			FinalPath: p,
			LockFile:  false,
		},
	)
}
