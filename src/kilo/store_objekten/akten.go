package store_objekten

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/id"
	"github.com/friedenberg/zit/src/delta/age_io"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
	"github.com/friedenberg/zit/src/india/zettel_verzeichnisse"
)

func (s Store) AkteExists(sh sha.Sha) (err error) {
	if sh.IsNull() {
		return
	}

	p := id.Path(sh, s.standort.DirObjektenAkten())
	ok := files.Exists(p)

	if !ok {
		return
	}

	set := zettel_transacted.MakeMutableSetUnique(0)

	w := zettel_verzeichnisse.MakeWriter(
		func(z *zettel_verzeichnisse.Zettel) (err error) {
			if !z.Transacted.Named.Stored.Zettel.Akte.Equals(sh) {
				err = io.EOF
				return
			}

			return
		},
	)

	if err = s.verzeichnisseAll.ReadMany(
		w,
		zettel_verzeichnisse.MakeWriterZettelTransacted(
			set.Add,
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	err = ErrAkteExists{
		Akte:       sh,
		MutableSet: set,
	}

	return
}

type akteMultiWriter struct {
	io.Writer
	writers []age_io.Writer
}

func (w akteMultiWriter) Close() (err error) {
	for _, w1 := range w.writers {
		if err = w1.Close(); err != nil {
			err = errors.Wrap(err)
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

func (s Store) AkteWriter() (w sha.WriteCloser, err error) {
	var outer age_io.Writer

	mo := age_io.MoveOptions{
		Age:                      s.age,
		FinalPath:                s.standort.DirObjektenAkten(),
		GenerateFinalPathFromSha: true,
		LockFile:                 true,
	}

	if outer, err = age_io.NewMover(mo); err != nil {
		err = errors.Wrap(err)
		return
	}

	w = outer

	return
}

func (s Store) AkteReader(sha sha.Sha) (r io.ReadCloser, err error) {
	if sha.IsNull() {
		//TODO move to files?
		r = ioutil.NopCloser(bytes.NewReader(nil))
		return
	}

	p := id.Path(sha, s.standort.DirObjektenAkten())

	o := age_io.FileReadOptions{
		Age:  s.age,
		Path: p,
	}

	if r, err = age_io.NewFileReader(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
