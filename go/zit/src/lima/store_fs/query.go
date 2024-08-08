package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

func (s *Store) MakeApplyCheckedOut(
	qg *query.Group,
	f interfaces.FuncIter[sku.CheckedOutLike],
	o sku.CommitOptions,
) interfaces.FuncIter[*FDSet] {
	return func(em *FDSet) (err error) {
		if err = s.ApplyCheckedOut(o, qg, em, f); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (s *Store) ApplyCheckedOut(
	o sku.CommitOptions,
	qg *query.Group,
	em *FDSet,
	f interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	var co *CheckedOut

	if co, err = s.ReadCheckedOutFromObjectIdFDPair(o, em); err != nil {
		err = errors.Wrapf(err, "%v", em)
		return
	}

	if co.External.FDs.State != external_state.Recognized &&
		!qg.ContainsExternalSku(
			&co.External.Transacted,
			co.State,
		) {
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

	o := sku.CommitOptions{
		Mode: objekte_mode.ModeRealizeSansProto,
	}

	aco := s.MakeApplyCheckedOut(qg, f, o)

	wg.Do(func() error {
		return s.AllObjects(aco)
	})

	if !qg.ExcludeUntracked {
		// wg.Do(func() error {
		// 	return s.QueryUnsure(qg, f)
		// })

		wg.Do(func() error {
			return s.QueryBlobs(qg, aco, f)
		})
	}

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryBlobs(
	qg *query.Group,
	aco interfaces.FuncIter[*FDSet],
	f func(sku.CheckedOutLike) error,
) (err error) {
	allRecognized := make([]*FDSet, 0)
	qg.SetIncludeHistory()

	if err = s.externalStoreInfo.FuncPrimitiveQuery(
		qg,
		func(sk *sku.Transacted) (err error) {
			shaBlob := sk.Metadata.Blob

			if shaBlob.IsNull() {
				return
			}

			key := shaBlob.GetBytes()
			recognized, ok := s.shasToBlobFDs[key]

			if !ok {
				return
			}

			recognizedFDS := &FDSet{
				State:          external_state.Recognized,
				MutableSetLike: collections_value.MakeMutableValueSet[*fd.FD](nil),
			}

			recognizedFDS.ObjectId.ResetWith(&sk.ObjectId)

			if err = recognized.Each(
				func(fds *FDSet) (err error) {
					fds.State = external_state.Recognized
					recognizedFDS.Add(fds.Blob.Clone())
					return
				},
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			allRecognized = append(allRecognized, recognizedFDS)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.dirFDs.blobs.Each(
		func(fds *FDSet) (err error) {
			if fds.State == external_state.Recognized {
				return
			}

			if err = aco(fds); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, fds := range allRecognized {
		if err = aco(fds); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
