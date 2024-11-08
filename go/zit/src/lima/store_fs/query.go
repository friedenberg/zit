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
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

func (s *Store) MakeApplyCheckedOut(
	qg *query.Group,
	f interfaces.FuncIter[sku.CheckedOutLike],
	o sku.CommitOptions,
) interfaces.FuncIter[*sku.FSItem] {
	return func(em *sku.FSItem) (err error) {
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
	i *sku.FSItem,
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
	aco interfaces.FuncIter[*sku.FSItem],
	f func(sku.CheckedOutLike) error,
) (err error) {
	allRecognizedBlobs := make([]*sku.FSItem, 0)
	allRecognizedObjects := make([]*sku.FSItem, 0)

	addRecognizedIfNecessary := func(
		sk *sku.Transacted,
		shaBlob *sha.Sha,
		shaCache map[sha.Bytes]interfaces.MutableSetLike[*sku.FSItem],
		allRecognized *[]*sku.FSItem, fdSetToFD func(*sku.FSItem) *fd.FD,
	) (err error) {
		if shaBlob.IsNull() {
			return
		}

		key := shaBlob.GetBytes()
		recognized, ok := shaCache[key]

		if !ok {
			return
		}

		recognizedFDS := &sku.FSItem{
			State:          external_state.Recognized,
			MutableSetLike: collections_value.MakeMutableValueSet[*fd.FD](nil),
		}

		recognizedFDS.ExternalObjectId.ResetWith(&sk.ObjectId)

		if err = recognized.Each(
			func(fds *sku.FSItem) (err error) {
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
				func(fds *sku.FSItem) *fd.FD { return &fds.Blob },
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = addRecognizedIfNecessary(
				sk,
				&sk.Metadata.SelfMetadataWithoutTai,
				s.shasToObjectFDs,
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

	if err = s.dirItems.ConsolidateDuplicateBlobs(); err != nil {
		err = errors.Wrap(err)
		return
	}

	{
		blobs := make([]*sku.FSItem, 0, s.dirItems.blobs.Len())

		if err = s.dirItems.blobs.Each(
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
		objects := make([]*sku.FSItem, 0, s.dirItems.objects.Len())

		if err = s.dirItems.objects.Each(
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
