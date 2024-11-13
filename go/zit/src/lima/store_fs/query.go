package store_fs

import (
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

func (s *Store) makeFuncIterHydrateCheckedOutProbablyCheckedOut(
	f interfaces.FuncIter[sku.SkuType],
) interfaces.FuncIter[*sku.FSItem] {
	return func(item *sku.FSItem) (err error) {
		var co *sku.CheckedOut

		if co, err = s.readCheckedOutFromItem(item); err != nil {
			err = errors.Wrapf(err, "%s", item.Debug())
			return
		}

		if err = s.WriteFSItemToExternal(item, co.GetSkuExternal()); err != nil {
			err = errors.Wrap(err)
			return
		}

		switch {
		case !item.Conflict.IsEmpty():
			co.SetState(checked_out_state.Conflicted)

			// case item.State == external_state.Recognized:
			// 	co.SetState(checked_out_state.Recognized)

			// case item.State == external_state.Untracked:
			// 	co.SetState(checked_out_state.Untracked)
		}

		if err = f(co); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (s *Store) makeFuncIterHydrateCheckedOutDefinitelyNotCheckedOut(
	f interfaces.FuncIter[sku.SkuType],
) interfaces.FuncIter[*sku.FSItem] {
	return func(item *sku.FSItem) (err error) {
		var co *sku.CheckedOut

		if co, err = s.readCheckedOutFromItem(item); err != nil {
			err = errors.Wrapf(err, "%s", item.Debug())
			return
		}

		if err = s.WriteFSItemToExternal(item, co.GetSkuExternal()); err != nil {
			err = errors.Wrap(err)
			return
		}

		switch {
		case !item.Conflict.IsEmpty():
			co.SetState(checked_out_state.Conflicted)

			// case item.State == external_state.Recognized:
			// 	co.SetState(checked_out_state.Recognized)

			// case item.State == external_state.Untracked:
			// 	co.SetState(checked_out_state.Untracked)
		}

		if err = f(co); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (s *Store) makeFuncIterFilterAndApply(
	qg *query.Group,
	f interfaces.FuncIter[sku.SkuType],
) interfaces.FuncIter[*sku.CheckedOut] {
	return func(co *sku.CheckedOut) (err error) {
		if !qg.ContainsExternalSku(
			co.GetSkuExternal(),
			co.GetState(),
		) {
			return
		}

		if err = f(co); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (s *Store) QueryCheckedOut(
	qg *query.Group,
	f interfaces.FuncIter[sku.SkuType],
) (err error) {
	wg := quiter.MakeErrorWaitGroupParallel()

	wg.Do(func() (err error) {
		aco := s.makeFuncIterHydrateCheckedOutProbablyCheckedOut(
			s.makeFuncIterFilterAndApply(qg, f),
		)

		for o := range s.probablyCheckedOut.All() {
			if err = aco(o); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		return
	})

	if !qg.ExcludeUntracked {
		wg.Do(func() (err error) {
			aco := s.makeFuncIterHydrateCheckedOutDefinitelyNotCheckedOut(
				s.makeFuncIterFilterAndApply(qg, f),
			)

			if err = s.queryUntracked(aco); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		})
	}

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) queryUntracked(
	aco interfaces.FuncIter[*sku.FSItem],
) (err error) {
	allRecognizedBlobs := make([]*sku.FSItem, 0)
	allRecognizedObjects := make([]*sku.FSItem, 0)

	addRecognizedIfNecessary := func(
		sk *sku.Transacted,
		shaBlob *sha.Sha,
		shaCache map[sha.Bytes]interfaces.MutableSetLike[*sku.FSItem],
		allRecognized *[]*sku.FSItem,
		fdSetToFD func(*sku.FSItem) *fd.FD,
	) (err error) {
		if shaBlob.IsNull() {
			return
		}

		key := shaBlob.GetBytes()
		recognized, ok := shaCache[key]

		if !ok {
			return
		}

		// TODO forward checked_out_state.Recognized
		recognizedFDS := &sku.FSItem{
			MutableSetLike: collections_value.MakeMutableValueSet[*fd.FD](nil),
		}

		recognizedFDS.ExternalObjectId.ResetWith(&sk.ObjectId)

		if err = recognized.Each(
			func(fds *sku.FSItem) (err error) {
				recognizedFDS.Add(fdSetToFD(fds).Clone())
				return
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		*allRecognized = append(*allRecognized, recognizedFDS)

		return
	}

	if err = s.externalStoreSupplies.FuncPrimitiveQuery(
		nil,
		func(sk *sku.Transacted) (err error) {
			if err = addRecognizedIfNecessary(
				sk,
				&sk.Metadata.Blob,
				s.definitelyNotCheckedOut.shas,
				&allRecognizedBlobs,
				func(fds *sku.FSItem) *fd.FD { return &fds.Blob },
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = addRecognizedIfNecessary(
				sk,
				&sk.Metadata.SelfMetadataWithoutTai,
				s.probablyCheckedOut.shas,
				&allRecognizedObjects,
				func(fds *sku.FSItem) *fd.FD { return &fds.Object },
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO move to initial parse?
	if err = s.dirItems.definitelyNotCheckedOut.ConsolidateDuplicateBlobs(); err != nil {
		err = errors.Wrap(err)
		return
	}

	{
		blobs := make([]*sku.FSItem, 0, s.dirItems.definitelyNotCheckedOut.Len())

		if err = s.dirItems.definitelyNotCheckedOut.Each(
			func(fds *sku.FSItem) (err error) {
				blobs = append(blobs, fds)
				return
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		sort.Slice(
			blobs,
			func(i, j int) bool {
				return blobs[i].ExternalObjectId.String() < blobs[j].ExternalObjectId.String()
			},
		)

		for _, fds := range blobs {
			// if fds.State == external_state.Recognized {
			// 	continue
			// }

			if err = aco(fds); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		for _, fds := range allRecognizedBlobs {
			if err = aco(fds); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	if false {
		objects := make([]*sku.FSItem, 0, s.dirItems.probablyCheckedOut.Len())

		if err = s.dirItems.probablyCheckedOut.Each(
			func(fds *sku.FSItem) (err error) {
				objects = append(objects, fds)
				return
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		sort.Slice(
			objects,
			func(i, j int) bool {
				return objects[i].ExternalObjectId.String() < objects[j].ExternalObjectId.String()
			},
		)

		for _, fds := range objects {
			// if fds.State == external_state.Recognized {
			// 	continue
			// }

			if err = aco(fds); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		for _, fds := range allRecognizedObjects {
			if err = aco(fds); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}
