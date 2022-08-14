package objekten

import (
	"io"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/charlie/sha"
	age_io "github.com/friedenberg/zit/echo/age_io"
	"github.com/friedenberg/zit/delta/id"
)

type akteMultiWriter struct {
	io.Writer
	writers []age_io.Writer
}

func (w akteMultiWriter) Close() (err error) {
	for _, w1 := range w.writers {
		if err = w1.Close(); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}

func (w akteMultiWriter) Sha() (s sha.Sha) {
	s = w.writers[0].Sha()

	for _, w1 := range w.writers[1:] {
		s1 := w1.Sha()
		if s1 != s {
			panic(
				errors.Errorf(
					"shas from multi-writer don't match:\nexpected: %s\nactual: %s\n",
					s,
					s1,
				),
			)
		}
	}

	return
}

func (s Store) AkteWriter() (w age_io.Writer, err error) {
	var outer age_io.Writer

	mo := age_io.MoveOptions{
		Age:                      s.Age,
		FinalPath:                s.Umwelt.DirObjektenAkten(),
		GenerateFinalPathFromSha: true,
	}

	if outer, err = age_io.NewMover(mo); err != nil {
		err = errors.Error(err)
		return
	}

	w = outer

	return
}

func (s Store) AkteReader(sha sha.Sha) (r io.ReadCloser, err error) {
	p := id.Path(sha, s.Umwelt.DirObjektenAkten())

	o := age_io.FileReadOptions{
		Age:  s.Age,
		Path: p,
	}

	if r, err = age_io.NewFileReader(o); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
