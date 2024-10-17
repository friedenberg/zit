package store_fs

import (
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

func (s *Store) MakeApplyCheckedOut(
	qg *query.Group,
	f interfaces.FuncIter[sku.CheckedOutLike],
	o sku.CommitOptions,
) interfaces.FuncIter[*Item] {
	return func(em *Item) (err error) {
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
	i *Item,
	f interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	var co *sku.CheckedOut

	if co, err = s.readCheckedOutFromItem(o, i); err != nil {
		err = errors.Wrapf(err, "%s", i.Debug())
		return
	}

	if err = s.WriteFSItemToExternal(i, &co.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !qg.ContainsExternalSku(
		&co.External,
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
	wg := quiter.MakeErrorWaitGroupParallel()

	o := sku.CommitOptions{
		Mode: object_mode.ModeRealizeSansProto,
	}

	aco := s.MakeApplyCheckedOut(qg, f, o)

	wg.Do(func() error {
		return s.OnlyObjects(aco)
	})

	if !qg.ExcludeUntracked {
		// wg.Do(func() error {
		// 	return s.QueryUnsure(qg, f)
		// })

		wg.Do(func() error {
			return s.QueryUntracked(qg, aco, f)
		})
	}

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryUntracked(
	qg *query.Group,
	aco interfaces.FuncIter[*Item],
	f func(sku.CheckedOutLike) error,
) (err error) {
	allRecognizedBlobs := make([]*Item, 0)
	allRecognizedObjects := make([]*Item, 0)

	addRecognizedIfNecessary := func(
		sk *sku.Transacted,
		shaBlob *sha.Sha,
		shaCache map[sha.Bytes]interfaces.MutableSetLike[*Item],
		allRecognized *[]*Item,
		fdSetToFD func(*Item) *fd.FD,
	) (err error) {
		if shaBlob.IsNull() {
			return
		}

		key := shaBlob.GetBytes()
		recognized, ok := shaCache[key]

		if !ok {
			return
		}

		recognizedFDS := &Item{
			State:          external_state.Recognized,
			MutableSetLike: collections_value.MakeMutableValueSet[*fd.FD](nil),
		}

		recognizedFDS.ObjectId.ResetWith(&sk.ObjectId)

		if err = recognized.Each(
			func(fds *Item) (err error) {
				fds.State = external_state.Recognized
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
				s.shasToBlobFDs,
				&allRecognizedBlobs,
				func(fds *Item) *fd.FD { return &fds.Blob },
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = addRecognizedIfNecessary(
				sk,
				&sk.Metadata.SelfMetadataWithoutTai,
				s.shasToObjectFDs,
				&allRecognizedObjects,
				func(fds *Item) *fd.FD { return &fds.Object },
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

	if err = s.dirItems.ConsolidateDuplicateBlobs(); err != nil {
		err = errors.Wrap(err)
		return
	}

	{
		blobs := make([]*Item, 0, s.dirItems.blobs.Len())

		if err = s.dirItems.blobs.Each(
			func(fds *Item) (err error) {
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
				return blobs[i].ObjectId.String() < blobs[j].ObjectId.String()
			},
		)

		for _, fds := range blobs {
			if fds.State == external_state.Recognized {
				continue
			}

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
		objects := make([]*Item, 0, s.dirItems.objects.Len())

		if err = s.dirItems.objects.Each(
			func(fds *Item) (err error) {
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
				return objects[i].ObjectId.String() < objects[j].ObjectId.String()
			},
		)

		for _, fds := range objects {
			if fds.State == external_state.Recognized {
				continue
			}

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
