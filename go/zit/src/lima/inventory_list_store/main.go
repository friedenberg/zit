package inventory_list_store

import (
	"iter"
	"sort"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_lock"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_immutable_io"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	"code.linenisgreat.com/zit/go/zit/src/kilo/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
)

type Store struct {
	lock sync.Mutex

	envRepo      env_repo.Env
	lockSmith    interfaces.LockSmith
	storeVersion interfaces.StoreVersion
	objectBlobStore
	blobStore blob_store.LocalBlobStore
	clock     ids.Clock

	object_format object_inventory_format.Format
	options       object_inventory_format.Options
	box           *box_format.BoxTransacted

	ui sku.UIStorePrinters
}

type objectBlobStore interface {
	getType() ids.Type
	getTypedBlobStore() typed_blob_store.InventoryList

	ReadOneSha(id interfaces.Stringer) (object *sku.Transacted, err error)
	WriteInventoryListObject(
		object *sku.Transacted,
	) (err error)

	IterAllInventoryLists() iter.Seq2[*sku.Transacted, error]
}

func (s *Store) Initialize(
	envRepo env_repo.Env,
	clock ids.Clock,
	typedBlobStore typed_blob_store.InventoryList,
) (err error) {
	op := object_inventory_format.Options{Tai: true}

	*s = Store{
		envRepo:      envRepo,
		lockSmith:    envRepo.GetLockSmith(),
		storeVersion: envRepo.GetStoreVersion(),
		blobStore:    envRepo.MakeBlobStore(),
		clock:        clock,
		box: box_format.MakeBoxTransactedArchive(
			envRepo,
			options_print.V0{}.WithPrintTai(true),
		),
		options: op,
	}

	v := s.storeVersion.GetInt()

	var blobType ids.Type

	switch {
	case v <= 6:
		blobType = ids.MustType(builtin_types.InventoryListTypeV0)

	default:
		blobType = ids.MustType(builtin_types.InventoryListTypeV1)
	}

	switch {
	case v <= 8:
		s.objectBlobStore = &objectBlobStoreV0{
			blobType: blobType,
			blobStore: blob_store.MakeShardedFilesStore(
				envRepo.DirInventoryLists(),
				env_dir.MakeConfigFromImmutableBlobConfig(
					envRepo.GetConfigPrivate().ImmutableConfig.GetBlobStoreConfigImmutable(),
				),
				envRepo.GetTempLocal(),
			),
			typedBlobStore: typedBlobStore,
		}

	default:
		s.objectBlobStore = &objectBlobStoreV1{
			pathLog:  envRepo.FileInventoryListLog(),
			blobType: blobType,
			blobStore: blob_store.MakeShardedFilesStore(
				envRepo.DirInventoryLists(),
				env_dir.MakeConfigFromImmutableBlobConfig(
					envRepo.GetConfigPrivate().ImmutableConfig.GetBlobStoreConfigImmutable(),
				),
				envRepo.GetTempLocal(),
			),
			typedBlobStore: typedBlobStore,
		}
	}

	return
}

func (s *Store) SetUIDelegate(ud sku.UIStorePrinters) {
	s.ui = ud
}

func (s *Store) GetEnv() env_ui.Env {
	return s.GetEnvRepo()
}

func (s *Store) GetImmutableConfigPublic() config_immutable_io.ConfigLoadedPublic {
	return s.GetEnvRepo().GetConfigPublic()
}

func (s *Store) GetImmutableConfigPrivate() config_immutable_io.ConfigLoadedPrivate {
	return s.GetEnvRepo().GetConfigPrivate()
}

func (s *Store) GetObjectStore() sku.ObjectStore {
	return s
}

func (s *Store) GetTypedInventoryListBlobStore() typed_blob_store.InventoryList {
	return s.getTypedBlobStore()
}

func (s *Store) Flush() (err error) {
	wg := errors.MakeWaitGroupParallel()
	return wg.GetError()
}

func (s *Store) FormatForVersion(sv interfaces.StoreVersion) sku.ListFormat {
	v := sv.GetInt()

	switch {
	case v <= 6:
		return inventory_list_blobs.MakeV0(
			s.object_format,
			s.options,
		)

	default:
		return inventory_list_blobs.V1{
			Box: s.box,
		}
	}
}

func (s *Store) GetTai() ids.Tai {
	if s.clock == nil {
		return ids.NowTai()
	} else {
		return s.clock.GetTai()
	}
}

func (s *Store) GetEnvRepo() env_repo.Env {
	return s.envRepo
}

func (s *Store) GetBlobStore() interfaces.BlobStore {
	return &s.envRepo
}

func (s *Store) GetInventoryListStore() sku.InventoryListStore {
	return s
}

func (store *Store) MakeOpenList() (openList *sku.OpenList, err error) {
	// TODO type
	openList = &sku.OpenList{}

	if openList.Mover, err = store.blobStore.Mover(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (store *Store) AddObjectToOpenList(
	openList *sku.OpenList,
	object *sku.Transacted,
) (err error) {
	format := store.FormatForVersion(store.storeVersion)

	if _, err = format.WriteObjectToOpenList(object, openList); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (store *Store) Create(
	openList *sku.OpenList,
) (object *sku.Transacted, err error) {
	if openList.Len == 0 {
		return
	}

	if !store.lockSmith.IsAcquired() {
		err = file_lock.ErrLockRequired{
			Operation: "create inventory list",
		}

		return
	}

	object = sku.GetTransactedPool().Get()

	object.Metadata.Type = store.getType()
	object.Metadata.Description = openList.Description

	tai := store.GetTai()

	if err = object.ObjectId.SetWithIdLike(tai); err != nil {
		err = errors.Wrap(err)
		return
	}

	object.SetTai(tai)

	if err = openList.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	actual := openList.GetShaLike()
	expected := sha.Make(object.GetBlobSha())

	ui.Log().Print("expected", expected, "actual", actual)

	if expected.IsNull() {
		object.SetBlobSha(actual)
	} else {
		if err = expected.AssertEqualsShaLike(actual); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = store.WriteInventoryListObject(object); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (store *Store) WriteInventoryListBlob(
	remoteBlobStore interfaces.BlobStore,
	object *sku.Transacted,
	list *sku.List,
) (err error) {
	if list.Len() == 0 {
		if !object.GetBlobSha().IsNull() {
			err = errors.ErrorWithStackf(
				"inventory list has non-empty blob but passed in list is empty. %q",
				sku.String(object),
			)

			return
		}

		return
	}

	var writeCloser interfaces.ShaWriteCloser

	if writeCloser, err = store.envRepo.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, writeCloser)

	if _, err = store.getTypedBlobStore().WriteBlobToWriter(
		object.GetType(),
		list,
		writeCloser,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	actual := writeCloser.GetShaLike()
	expected := sha.Make(object.GetBlobSha())

	ui.Log().Print("expected", expected, "actual", actual)

	if expected.IsNull() {
		object.SetBlobSha(actual)
	} else {
		if err = expected.AssertEqualsShaLike(actual); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	// if !s.af.HasBlob(t.GetBlobSha()) {
	// 	err = errors.Errorf(
	// 		"inventory list blob missing after write (%d bytes, %d skus): %q",
	// 		n,
	// 		skus.Len(),
	// 		sku.String(t),
	// 	)

	// 	return
	// }

	// if _, _, err = s.blobStore.GetTransactedWithBlob(
	// 	t,
	// ); err != nil {
	// 	err = errors.Wrapf(err, "Blob Sha: %q", actual)
	// 	return
	// }

	return
}

// TODO split into public and private parts, where public includes writing the
// skus AND the list, while private writes just the list
func (store *Store) ImportInventoryList(
	remoteBlobStore interfaces.BlobStore,
	object *sku.Transacted,
) (err error) {
	var blobReader interfaces.ShaReadCloser

	if blobReader, err = remoteBlobStore.BlobReader(
		object.GetBlobSha(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	list := sku.MakeList()

	if err = inventory_list_blobs.ReadInventoryListBlob(
		store.FormatForVersion(store.storeVersion),
		blobReader,
		list,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	for sk := range list.All() {
		if err = sk.CalculateObjectShas(); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = blob_store.CopyBlobIfNecessary(
			store.GetEnvRepo().GetEnv(),
			store.blobStore,
			remoteBlobStore,
			sk.GetBlobSha(),
			nil,
		); err != nil {
			if errors.Is(err, &env_dir.ErrAlreadyExists{}) {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}

			continue
		}
	}

	if err = store.WriteInventoryListBlob(
		remoteBlobStore,
		object,
		list,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = store.WriteInventoryListObject(
		object,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) IterInventoryList(
	blobSha interfaces.Sha,
) iter.Seq2[*sku.Transacted, error] {
	return s.getTypedBlobStore().IterInventoryListBlobSkusFromBlobStore(
		s.getType(),
		s.blobStore,
		blobSha,
	)
}

func (store *Store) ReadLast() (max *sku.Transacted, err error) {
	var maxSku sku.Transacted

	for b, iterErr := range store.IterAllInventoryLists() {
		if iterErr != nil {
			err = errors.Wrap(iterErr)
			return
		}

		if sku.TransactedLessor.LessPtr(&maxSku, b) {
			sku.TransactedResetter.ResetWith(&maxSku, b)
		}
	}

	max = &maxSku

	return
}

func (store *Store) ReadAllSorted(
	output interfaces.FuncIter[*sku.Transacted],
) (err error) {
	var skus []*sku.Transacted

	for list, iterErr := range store.IterAllInventoryLists() {
		if iterErr != nil {
			err = errors.Wrap(iterErr)
			return
		}

		skus = append(skus, list)
	}

	sort.Slice(skus, func(i, j int) bool { return skus[i].Less(skus[j]) })

	for _, o := range skus {
		if err = output(o); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (store *Store) IterAllSkus() iter.Seq2[sku.ObjectWithList, error] {
	return func(yield func(sku.ObjectWithList, error) bool) {
		var objectWithList sku.ObjectWithList

		for listObject, iterErr := range store.IterAllInventoryLists() {
			objectWithList.List = listObject
			objectWithList.Object = listObject

			if iterErr != nil {
				if !yield(objectWithList, iterErr) {
					return
				}
			}

			if !yield(objectWithList, nil) {
				return
			}

			iter := store.IterInventoryList(
				listObject.GetBlobSha(),
			)

			for object, iterErr := range iter {
				objectWithList.Object = object

				if !yield(objectWithList, iterErr) {
					return
				}
			}
		}
	}
}

func (store *Store) ReadAllSkus(
	f func(listSku, sk *sku.Transacted) error,
) (err error) {
	for listObject, iterErr := range store.IterAllInventoryLists() {
		if iterErr != nil {
			err = errors.Wrap(iterErr)
			return
		}

		if err = f(listObject, listObject); err != nil {
			err = errors.Wrapf(
				err,
				"InventoryList: %s",
				listObject.GetObjectId(),
			)

			return
		}

		iter := store.IterInventoryList(
			listObject.GetBlobSha(),
		)

		for object, iterErr := range iter {
			if iterErr != nil {
				if object == nil {
					err = errors.Wrap(iterErr)
				} else {
					err = errors.Wrapf(iterErr, "Sku: %s", sku.String(object))
				}

				return
			}

			if err = f(listObject, object); err != nil {
				err = errors.Wrapf(
					err,
					"InventoryList: %s",
					listObject.GetObjectId(),
				)

				return
			}
		}
	}

	return
}
