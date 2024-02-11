package store_objekten

import (
	"sync"

	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/bravo/files"
	"code.linenisgreat.com/zit-go/src/bravo/id"
	"code.linenisgreat.com/zit-go/src/bravo/iter"
	"code.linenisgreat.com/zit-go/src/charlie/gattung"
	"code.linenisgreat.com/zit-go/src/delta/gattungen"
	"code.linenisgreat.com/zit-go/src/echo/fd"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
)

func (s Store) ReadAllMatchingAkten(
	akten fd.Set,
	f func(*fd.FD, *sku.Transacted) error,
) (err error) {
	fds := fd.MakeMutableSetSha()

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
			func(fd *fd.FD) (err error) {
				if fd.GetShaLike().IsNull() {
					return iter.MakeErrStopIteration()
				}

				p := id.Path(fd.GetShaLike(), pa)

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

	observed := fd.MakeMutableSet()
	var l sync.Mutex

	if err = s.ReadAll(
		gattungen.MakeSet(gattung.TrueGattung()...),
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
		func(fd *fd.FD) (err error) {
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
