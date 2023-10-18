package store_objekten

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type KeyerFDSha struct{}

func (k KeyerFDSha) GetKey(fd fd.FD) string {
	return fd.Sha.String()
}

func (s Store) ReadAllMatchingAkten(
	akten schnittstellen.SetLike[fd.FD],
	f func(fd.FD, *sku.Transacted) error,
) (err error) {
	fds := collections_value.MakeMutableValueSet[fd.FD](
		KeyerFDSha{},
	)

	var pa string

	if pa, err = s.GetStandort().DirObjektenGattung(
		s.GetKonfig().GetStoreVersion(),
		gattung.Akte,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = akten.Each(
		iter.MakeChain(
			func(fd fd.FD) (err error) {
				if fd.Sha.IsNull() {
					return iter.MakeErrStopIteration()
				}

				p := id.Path(fd.Sha, pa)

				if !files.Exists(p) {
					return iter.MakeErrStopIteration()
				}

				return
			},
			// TODO-P2 handle files with the same sha
			fds.Add,
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	observed := collections_value.MakeMutableValueSet[fd.FD](nil)
	var l sync.Mutex

	if err = s.Zettel().ReadAll(
		func(z *sku.Transacted) (err error) {
			fd, ok := fds.Get(z.GetAkteSha().String())

			if !ok {
				return
			}

			if err = f(fd, z); err != nil {
				err = errors.Wrap(err)
				return
			}

			l.Lock()
			defer l.Unlock()

			return observed.Add(fd)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = fds.Each(
		func(fd fd.FD) (err error) {
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
