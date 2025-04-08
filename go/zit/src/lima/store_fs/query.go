package store_fs

import (
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

func (s *Store) QueryCheckedOut(
	qg *query.Query,
	f interfaces.FuncIter[sku.SkuType],
) (err error) {
	wg := errors.MakeWaitGroupParallel()

	wg.Do(func() (err error) {
		funcIterFSItems := s.makeFuncIterHydrateCheckedOutProbablyCheckedOut(
			s.makeFuncIterFilterAndApply(qg, f),
		)

		for o := range s.probablyCheckedOut.All() {
			if err = funcIterFSItems(o); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		return
	})

	if !qg.ExcludeUntracked {
		wg.Do(func() (err error) {
			funcIterFSItems := s.makeFuncIterHydrateCheckedOutDefinitelyNotCheckedOut(
				s.makeFuncIterFilterAndApply(qg, f),
			)

			if err = s.queryUntracked(qg, funcIterFSItems); err != nil {
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

func (s *Store) makeFuncIterHydrateCheckedOutProbablyCheckedOut(
	out interfaces.FuncIter[sku.SkuType],
) interfaces.FuncIter[*sku.FSItem] {
	return func(item *sku.FSItem) (err error) {
		co := GetCheckedOutPool().Get()

		// at a bare minimum, the internal object ID must always be set as there are
		// hard assumptions about internal being valid throughout the reading cycle
		if err = co.GetSku().ObjectId.SetObjectIdLike(&item.ExternalObjectId); err != nil {
			err = errors.Wrap(err)
			return
		}

		hasInternal := true

		var oid ids.ObjectId

		if err = oid.SetObjectIdLike(item.GetExternalObjectId()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = s.storeSupplies.ReadOneInto(
			&oid,
			co.GetSku(),
		); err != nil {
			if collections.IsErrNotFound(err) || genres.IsErrUnsupportedGenre(err) {
				hasInternal = false
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		if err = s.HydrateExternalFromItem(
			sku.CommitOptions{
				StoreOptions: sku.StoreOptions{
					UpdateTai: true,
				},
			},
			item,
			co.GetSku(),
			co.GetSkuExternal(),
		); err != nil {
			if sku.IsErrMergeConflict(err) {
				co.SetState(checked_out_state.Conflicted)

				if err = co.GetSkuExternal().ObjectId.SetWithIdLike(
					&co.GetSku().ObjectId,
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			} else {
				err = errors.Wrapf(err, "Cwd: %#v", item.Debug())
				return
			}
		}

		if !item.Conflict.IsEmpty() {
			co.SetState(checked_out_state.Conflicted)
		} else if !hasInternal {
			co.SetState(checked_out_state.Untracked)
		} else {
			co.SetState(checked_out_state.CheckedOut)
		}

		if err = s.WriteFSItemToExternal(item, co.GetSkuExternal()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = out(co); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (s *Store) makeFuncIterHydrateCheckedOutDefinitelyNotCheckedOut(
	f interfaces.FuncIter[sku.SkuType],
) interfaces.FuncIter[any] {
	return func(itemUnknown any) (err error) {
		co := sku.GetCheckedOutPool().Get()

		switch item := itemUnknown.(type) {
		case *sku.FSItem:
			if err = s.hydrateDefinitelyNotCheckedOutUnrecognizedItem(
				item,
				co,
				f,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

		case *fsItemRecognized:
			if err = s.hydrateDefinitelyNotCheckedOutRecognizedItem(
				item,
				co,
				f,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

		default:
			err = errors.Errorf("unsupported type for item: %T", itemUnknown)
			return
		}

		return
	}
}

func (s *Store) hydrateDefinitelyNotCheckedOutUnrecognizedItem(
	item *sku.FSItem,
	co *sku.CheckedOut,
	f interfaces.FuncIter[sku.SkuType],
) (err error) {
	if !item.Conflict.IsEmpty() {
		err = errors.Errorf("cannot have a conflict for a definitely not checked out blob: %s", item.Debug())
		return
	}

	if err = co.GetSku().ObjectId.SetObjectIdLike(
		&item.ExternalObjectId,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = co.GetSkuExternal().ObjectId.SetObjectIdLike(
		&item.ExternalObjectId,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.readOneExternalBlob(
		co.GetSkuExternal(),
		co.GetSku(),
		item,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.WriteFSItemToExternal(item, co.GetSkuExternal()); err != nil {
		err = errors.Wrap(err)
		return
	}

	co.SetState(checked_out_state.Untracked)

	if err = f(co); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) hydrateDefinitelyNotCheckedOutRecognizedItem(
	item *fsItemRecognized,
	co *sku.CheckedOut,
	f interfaces.FuncIter[sku.SkuType],
) (err error) {
	sku.TransactedResetter.ResetWith(co.GetSku(), &item.Recognized)
	sku.TransactedResetter.ResetWith(co.GetSkuExternal(), &item.Recognized)

	co.SetState(checked_out_state.Recognized)

	for _, item := range item.Matching {
		if err = item.WriteToSku(
			co.GetSkuExternal(),
			s.envRepo,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		co.GetSkuExternal().ObjectId.SetGenre(genres.Blob)

		if err = s.WriteFSItemToExternal(item, co.GetSkuExternal()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = f(co); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) makeFuncIterFilterAndApply(
	qg *query.Query,
	f interfaces.FuncIter[sku.SkuType],
) interfaces.FuncIter[*sku.CheckedOut] {
	return func(co *sku.CheckedOut) (err error) {
		if !query.ContainsExternalSku(
			qg,
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

type fsItemRecognized struct {
	Recognized sku.Transacted
	Matching   []*sku.FSItem
}

func (s *Store) queryUntracked(
	qg *query.Query, // TODO use this to conditionally perform recognition
	aco interfaces.FuncIter[any],
) (err error) {
	definitelyNotCheckedOut := s.dirInfo.definitelyNotCheckedOut.Clone()

	// TODO move to initial parse?
	if err = definitelyNotCheckedOut.ConsolidateDuplicateBlobs(); err != nil {
		err = errors.Wrap(err)
		return
	}

	allRecognized := make([]*fsItemRecognized, 0)

	addRecognizedIfNecessary := func(
		sk *sku.Transacted,
		shaBlob *sha.Sha,
		shaCache map[sha.Bytes]interfaces.MutableSetLike[*sku.FSItem],
	) (item *fsItemRecognized, err error) {
		if shaBlob.IsNull() {
			return
		}

		key := shaBlob.GetBytes()
		recognized, ok := shaCache[key]

		if !ok {
			return
		}

		item = &fsItemRecognized{}

		sku.TransactedResetter.ResetWith(&item.Recognized, sk)

		for recognized := range recognized.All() {
			item.Matching = append(item.Matching, recognized)
		}

		return
	}

	if err = s.storeSupplies.ReadPrimitiveQuery(
		nil,
		func(sk *sku.Transacted) (err error) {
			var recognizedBlob, recognizedObject *fsItemRecognized

			if recognizedBlob, err = addRecognizedIfNecessary(
				sk,
				&sk.Metadata.Blob,
				definitelyNotCheckedOut.shas,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if recognizedObject, err = addRecognizedIfNecessary(
				sk,
				&sk.Metadata.SelfMetadataWithoutTai,
				s.probablyCheckedOut.shas,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if recognizedBlob != nil {
				allRecognized = append(allRecognized, recognizedBlob)

				for _, item := range recognizedBlob.Matching {
					definitelyNotCheckedOut.Del(item)
				}
			}

			if recognizedObject != nil {
				allRecognized = append(allRecognized, recognizedObject)

				for _, item := range recognizedObject.Matching {
					definitelyNotCheckedOut.Del(item)
				}
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	{
		blobs := make([]*sku.FSItem, 0, definitelyNotCheckedOut.Len())

		if err = definitelyNotCheckedOut.Each(
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

	}

	// if false {
	// 	objects := make([]*sku.FSItem, 0, s.dirItems.probablyCheckedOut.Len())

	// 	if err = s.dirItems.probablyCheckedOut.Each(
	// 		func(fds *sku.FSItem) (err error) {
	// 			objects = append(objects, fds)
	// 			return
	// 		},
	// 	); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// 	sort.Slice(
	// 		objects,
	// 		func(i, j int) bool {
	// 			return objects[i].ExternalObjectId.String() < objects[j].ExternalObjectId.String()
	// 		},
	// 	)

	// 	for _, fds := range objects {
	// 		// if fds.State == external_state.Recognized {
	// 		// 	continue
	// 		// }

	// 		if err = aco(fds); err != nil {
	// 			err = errors.Wrap(err)
	// 			return
	// 		}
	// 	}
	// }

	for _, fds := range allRecognized {
		if err = aco(fds); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
