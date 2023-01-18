package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/foxtrot/id"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

// TODO-P2 move to Standort
func (s Store) ReadAllAktenShas(w collections.WriterFunc[sha.Sha]) (err error) {
	wf := func(p string) (err error) {
		var sh sha.Sha

		if sh, err = sha.MakeShaFromPath(p); err != nil {
			err = errors.Wrapf(err, "Path: %s", p)
			return
		}

		if err = w(sh); err != nil {
			err = errors.Wrapf(err, "Sha: %s", sh)
			return
		}

		return
	}

	if err = files.ReadDirNamesLevel2(
		files.MakeDirNameWriterIgnoringHidden(wf),
		s.common.Standort.DirObjektenAkten(),
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

	p := id.Path(sh, s.common.Standort.DirObjektenAkten())
	ok := files.Exists(p)

	if !ok {
		return
	}

	set := zettel.MakeMutableSetUnique(0)

	if err = s.zettelStore.verzeichnisseAll.ReadMany(
		func(z *zettel.Transacted) (err error) {
			if !z.Objekte.Akte.Equals(sh) {
				err = collections.ErrStopIteration
				return
			}

			return
		},
		set.AddAndDoNotRepool,
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
