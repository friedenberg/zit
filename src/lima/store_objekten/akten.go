package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

func (s Store) ReadAllMatchingAkten(
	akten schnittstellen.Set[kennung.FD],
	f func(kennung.FD, *zettel.Transacted) error,
) (err error) {
	fds := collections.MakeMutableSet[kennung.FD](
		func(fd kennung.FD) string {
			return fd.Sha.String()
		},
	)

	if err = akten.Each(
		iter.MakeChain(
			func(fd kennung.FD) (err error) {
				if fd.Sha.IsNull() {
					return iter.MakeErrStopIteration()
				}

				p := id.Path(fd.Sha, s.StoreUtil.GetStandort().DirObjektenAkten())

				if !files.Exists(p) {
					return iter.MakeErrStopIteration()
				}

				return
			},
			//TODO handle files with the same sha
			fds.Add,
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	observed := collections.MakeMutableSetStringer[kennung.FD]()

	if err = s.Zettel().ReadAll(
		func(z *zettel.Transacted) (err error) {
			fd, ok := fds.Get(z.Objekte.Akte.String())

			if !ok {
				return
			}

			if err = f(fd, z); err != nil {
				err = errors.Wrap(err)
				return
			}

			return observed.Add(fd)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = fds.Each(
		func(fd kennung.FD) (err error) {
			if observed.Contains(fd) {
				return
			}

			if err = f(fd, nil); err != nil {
				err = errors.Wrap(err)
				return
			}

			return observed.Add(fd)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
