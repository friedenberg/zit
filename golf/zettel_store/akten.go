package zettel_store

import (
	"io"
	"os"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/delta/objekte"
)

type akteMultiWriter struct {
	io.Writer
	writers []objekte.Writer
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

func (s ZettelStore) AkteWriter() (w objekte.Writer, err error) {
	var inner, outer objekte.Writer

	if inner, err = s.Zettels.AkteWriter(); err != nil {
		err = errors.Error(err)
		return
	}

	if outer, err = objekte.NewWriterMover(s.Age(), s.Umwelt().DirAkten()); err != nil {
		err = errors.Error(err)
		return
	}

	w = akteMultiWriter{
		Writer:  io.MultiWriter(inner, outer),
		writers: []objekte.Writer{inner, outer},
	}

	return
}

type akteReader struct {
	file *os.File
	objekte.Reader
}

func (ar akteReader) Close() (err error) {
	if ar.file == nil {
		err = errors.Errorf("nil file")
		return
	}

	if ar.Reader == nil {
		err = errors.Errorf("nil objekte reader")
		return
	}

	if err = ar.Reader.Close(); err != nil {
		err = errors.Error(err)
		return
	}

	// if err = open_file_guard.Close(ar.file); err != nil {
	// 	err = errors.Error(err)
	// 	return
	// }

	return
}

func (s ZettelStore) AkteReader(sha sha.Sha) (r io.ReadCloser, err error) {
	ar := akteReader{}

	p := id.Path(sha, s.Umwelt().DirAkten())

	if ar.file, err = open_file_guard.Open(p); err != nil {
		if os.IsNotExist(err) {
			return s.Zettels.AkteReader(sha)
		} else {
			err = errors.Error(err)
			return
		}
	}

	if ar.Reader, err = objekte.NewReader(s.Age(), ar.file); err != nil {
		err = errors.Error(err)
		return
	}

	r = ar

	return
}
