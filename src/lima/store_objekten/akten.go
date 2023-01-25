package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

func (s Store) AkteExists(sh sha.Sha) (err error) {
  errors.TodoP3("decide what to do with this method")

	if sh.IsNull() {
		return
	}

	p := id.Path(sh, s.common.GetStandort().DirObjektenAkten())
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
