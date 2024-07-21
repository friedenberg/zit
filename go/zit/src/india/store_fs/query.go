package store_fs

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

func (s *Store) MakeApplyCheckedOut(
	qg sku.Queryable,
	f interfaces.FuncIter[sku.CheckedOutLike],
	o sku.CommitOptions,
) interfaces.FuncIter[*ObjectIdFDPair] {
	return func(em *ObjectIdFDPair) (err error) {
		if err = s.ApplyCheckedOut(o, qg, em, f); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (s *Store) ApplyCheckedOut(
	o sku.CommitOptions,
	qg sku.Queryable,
	em *ObjectIdFDPair,
	f interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	var co *CheckedOut

	if co, err = s.ReadCheckedOutFromObjectIdFDPair(o, em); err != nil {
		err = errors.Wrapf(err, "%v", em)
		return
	}

	// ui.Debug().Print(qg, qg.ContainsSku(&co.External.Transacted), co)

	if !qg.ContainsSku(&co.External.Transacted) {
		return
	}

	if err = f(co); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryCheckedOut(
	qg *query.Group,
	f interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	wg := iter.MakeErrorWaitGroupParallel()

	{
		o := sku.CommitOptions{
			Mode: objekte_mode.ModeRealizeSansProto,
		}

		wg.Do(func() error {
			return s.All(s.MakeApplyCheckedOut(qg, f, o))
		})
	}

	if !qg.ExcludeUntracked {
		wg.Do(func() error {
			return s.QueryUnsure(qg, f)
		})

		wg.Do(func() error {
			return s.QueryUntrackedBlobs(qg, f)
		})
	}

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryUnsure(
	qg *query.Group,
	f interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	o := sku.CommitOptions{
		Mode: objekte_mode.ModeRealizeWithProto,
	}

	if err = s.AllUnsure(
		s.MakeApplyCheckedOut(qg, f, o),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryUntrackedBlobs(
	qg *query.Group,
	f interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	if err = s.QueryAllMatchingBlobs(
		qg,
		s.GetUnsureBlobs(),
		func(fd *fd.FD, z *sku.Transacted) (err error) {
			fr := GetCheckedOutPool().Get()
			defer GetCheckedOutPool().Put(fr)

			fr.External.FDs.Blob.ResetWith(fd)
			fr.External.Metadata.Tai = ids.TaiFromTime(fd.ModTime())

			if z == nil {
				// TODO use ReadOneExternalBlob
				fr.State = checked_out_state.StateUntracked
				fr.External.SetBlobSha(fd.GetShaLike())

				if err = fr.External.Transacted.CalculateObjectShas(); err != nil {
					err = errors.Wrap(err)
					return
				}
			} else {
				fr.External.SetBlobSha(z.GetBlobSha())
				fr.State = checked_out_state.StateRecognized

				if err = fr.Internal.SetFromSkuLike(z); err != nil {
					err = errors.Wrap(err)
					return
				}

				sku.Resetter.ResetWith(&fr.External, z)

				if err = fr.External.SetObjectSha(z.GetObjectSha()); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			if err = f(fr); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryAllMatchingBlobs(
	qg *query.Group,
	blob_store fd.Set,
	f func(*fd.FD, *sku.Transacted) error,
) (err error) {
	fds := fd.MakeMutableSetSha()

	var pa string

	if pa, err = s.externalStoreInfo.Home.DirObjektenGattung(
		s.config.GetStoreVersion(),
		genres.Blob,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = blob_store.Each(
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

	if err = s.externalStoreInfo.FuncPrimitiveQuery(
		qg,
		func(z *sku.Transacted) (err error) {
			fd, ok := fds.Get(z.GetBlobSha().String())

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
