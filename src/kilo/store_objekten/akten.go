package store_objekten

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/id"
	"github.com/friedenberg/zit/src/delta/age_io"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
	"github.com/friedenberg/zit/src/india/zettel_verzeichnisse"
)

func (s Store) ReadAllAktenShas(w collections.WriterFunc[sha.Sha]) (err error) {
	wf := func(p string) (err error) {
		var sh sha.Sha

		if sh, err = sha.MakeShaFromPath(p); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = w(sh); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = files.ReadDirNamesLevel2(
		files.MakeDirNameWriterIgnoringHidden(wf),
		s.standort.DirObjektenAkten(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

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

	if err = s.verzeichnisseAll.ReadMany(
		func(z *zettel_verzeichnisse.Zettel) (err error) {
			if !z.Transacted.Named.Stored.Zettel.Akte.Equals(sh) {
				err = io.EOF
				return
			}

			return
		},
		zettel_verzeichnisse.MakeWriterZettelTransacted(
			set.AddAndDoNotRepool,
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
