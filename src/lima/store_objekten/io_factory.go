package store_objekten

import (
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/age_io"
)

func (s Store) ReadCloserObjekten(p string) (sha.ReadCloser, error) {
	o := age_io.FileReadOptions{
		Age:  s.age,
		Path: p,
	}

	return age_io.NewFileReader(o)
}

func (s Store) ReadCloserVerzeichnisse(p string) (sha.ReadCloser, error) {
	o := age_io.FileReadOptions{
		Age:  s.age,
		Path: p,
	}

	return age_io.NewFileReader(o)
}

func (s Store) WriteCloserObjekten(p string) (w sha.WriteCloser, err error) {
	return age_io.NewMover(
		age_io.MoveOptions{
			Age:       s.age,
			FinalPath: p,
			LockFile:  true,
		},
	)
}

func (s Store) WriteCloserVerzeichnisse(p string) (w sha.WriteCloser, err error) {
	return age_io.NewMover(
		age_io.MoveOptions{
			Age:       s.age,
			FinalPath: p,
			LockFile:  false,
		},
	)
}
