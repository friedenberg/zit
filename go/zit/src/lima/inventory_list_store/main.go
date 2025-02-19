package inventory_list_store

import (
	"iter"
	"sort"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_lock"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_immutable_io"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	"code.linenisgreat.com/zit/go/zit/src/kilo/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
)

type Store struct {
	lock sync.Mutex

	envRepo         env_repo.Env
	lockSmith       interfaces.LockSmith
	storeVersion    interfaces.StoreVersion
	objectBlobStore interfaces.BlobStore
	blobStore       interfaces.BlobStore
	clock           ids.Clock
	typedBlobStore  typed_blob_store.InventoryList

	object_format object_inventory_format.Format
	options       object_inventory_format.Options
	box           *box_format.BoxTransacted

	blobType ids.Type

	ui sku.UIStorePrinters
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
		objectBlobStore: env_repo.MakeBlobStore(
			envRepo.DirInventoryLists(),
			env_dir.Config{},
			envRepo.GetTempLocal(),
		),
		blobStore: envRepo,
		clock:     clock,
		box: box_format.MakeBoxTransactedArchive(
			envRepo,
			options_print.V0{}.WithPrintTai(true),
		),
		options:        op,
		typedBlobStore: typedBlobStore,
	}

	v := s.storeVersion.GetInt()

	switch {
	case v <= 6:
		s.blobType = ids.MustType(builtin_types.InventoryListTypeV0)

	default:
		s.blobType = ids.MustType(builtin_types.InventoryListTypeV1)
	}

	return
}

func (s *Store) SetUIDelegate(ud sku.UIStorePrinters) {
	s.ui = ud
}

func (s *Store) GetEnv() env_ui.Env {
	return s.GetEnvRepo()
}

func (s *Store) GetImmutableConfig() config_immutable_io.ConfigLoaded {
	return s.GetEnvRepo().GetConfig()
}

func (s *Store) GetObjectStore() sku.ObjectStore {
	return s
}

func (s *Store) GetTypedInventoryListBlobStore() typed_blob_store.InventoryList {
	return s.typedBlobStore
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

func (store *Store) Create(
	skus *sku.List,
	description descriptions.Description,
) (t *sku.Transacted, err error) {
	if skus.Len() == 0 {
		return
	}

	if !store.lockSmith.IsAcquired() {
		err = file_lock.ErrLockRequired{
			Operation: "create inventory list",
		}

		return
	}

	t = sku.GetTransactedPool().Get()

	t.Metadata.Type = store.blobType
	t.Metadata.Description = description

	tai := store.GetTai()

	if err = t.ObjectId.SetWithIdLike(tai); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.SetTai(tai)

	if err = store.WriteInventoryListBlob(t, skus); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = store.WriteInventoryListObject(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) WriteInventoryListBlob(
	t *sku.Transacted,
	skus *sku.List,
) (err error) {
	if skus.Len() == 0 {
		if !t.GetBlobSha().IsNull() {
			err = errors.Errorf(
				"inventory list has non-empty blob but passed in list is empty. %q",
				sku.String(t),
			)

			return
		}

		return
	}

	var wc interfaces.ShaWriteCloser

	if wc, err = s.envRepo.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, wc)

	if _, err = s.typedBlobStore.WriteBlobToWriter(
		t.GetType(),
		skus,
		wc,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	actual := wc.GetShaLike()
	expected := sha.Make(t.GetBlobSha())

	ui.Log().Print("expected", expected, "actual", actual)

	if expected.IsNull() {
		t.SetBlobSha(actual)
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
func (store *Store) WriteInventoryListObject(
	object *sku.Transacted,
) (err error) {
	store.lock.Lock()
	defer store.lock.Unlock()

	var wc interfaces.ShaWriteCloser

	// TODO also write to inventory_list_log

	if wc, err = store.objectBlobStore.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, wc)

	object.Metadata.Type = store.blobType

	if _, err = store.typedBlobStore.WriteObjectToWriter(
		object,
		wc,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = object.CalculateObjectShas(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ui.Log().Printf(
		"saved inventory list: %q",
		sku.String(object),
	)

	return
}

func (s *Store) ImportInventoryList(
	bs interfaces.BlobStore,
	t *sku.Transacted,
) (err error) {
	var rc interfaces.ShaReadCloser

	if rc, err = bs.BlobReader(
		t.GetBlobSha(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	list := sku.MakeList()

	if err = inventory_list_blobs.ReadInventoryListBlob(
		s.FormatForVersion(s.storeVersion),
		rc,
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

		if _, err = env_repo.CopyBlobIfNecessary(
			s.GetEnvRepo().GetEnv(),
			s.GetEnvRepo(),
			bs,
			sk.GetBlobSha(),
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

	if err = s.WriteInventoryListBlob(
		t,
		list,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.WriteInventoryListObject(
		t,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) readOnePath(p string) (o *sku.Transacted, err error) {
	var sh *sha.Sha

	if sh, err = sha.MakeShaFromPath(p); err != nil {
		err = errors.Wrapf(err, "Path: %q", p)
		return
	}

	if o, err = s.ReadOneSha(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = o.CalculateObjectShas(); err != nil {
		if errors.Is(err, object_inventory_format.ErrEmptyTai) {
			var t ids.Tai
			err1 := t.Set(o.ObjectId.String())

			if err1 != nil {
				err = errors.Wrapf(err, "%s", sku.StringTaiGenreObjectIdShaBlob(o))
				return
			}

			o.SetTai(t)

			if err = o.CalculateObjectShas(); err != nil {
				err = errors.Wrapf(err, "%#v", o)
				return
			}
		} else {
			err = errors.Wrapf(err, "%#v", o)
		}

		return
	}

	return
}

func (s *Store) ReadOneSha(
	k interfaces.Stringer,
) (o *sku.Transacted, err error) {
	var sh sha.Sha

	if err = sh.Set(k.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var or sha.ReadCloser

	if or, err = s.objectBlobStore.BlobReader(&sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, or)

	if o, err = s.typedBlobStore.ReadInventoryListObject(
		s.blobType,
		or,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) IterInventoryList(
	blobSha interfaces.Sha,
) iter.Seq2[*sku.Transacted, error] {
	return s.typedBlobStore.IterInventoryListBlobSkusFromBlobStore(
		s.blobType,
		s.blobStore,
		blobSha,
	)
}

func (s *Store) ReadLast() (max *sku.Transacted, err error) {
	var maxSku sku.Transacted

	for b, iterErr := range s.IterAllInventoryLists() {
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

// TODO switch to using append-only log
func (s *Store) IterAllInventoryLists() iter.Seq2[*sku.Transacted, error] {
	return func(yield func(*sku.Transacted, error) bool) {
		dir := s.envRepo.DirInventoryLists()

		for path, err := range files.DirNamesLevel2(dir) {
			if err != nil {
				if !yield(nil, errors.Wrap(err)) {
					return
				}
			}

			var decodedList *sku.Transacted

			{
				var err error

				if decodedList, err = s.readOnePath(path); err != nil {
					if !yield(nil, errors.Wrap(err)) {
						return
					}
				}
			}

			if !yield(decodedList, nil) {
				return
			}
		}
	}
}

func (s *Store) ReadAllSorted(
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	var skus []*sku.Transacted

	for list, iterErr := range s.IterAllInventoryLists() {
		if iterErr != nil {
			err = errors.Wrap(iterErr)
			return
		}

		skus = append(skus, list)
	}

	sort.Slice(skus, func(i, j int) bool { return skus[i].Less(skus[j]) })

	for _, o := range skus {
		if err = f(o); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (store *Store) ReadAllSkus(
	f func(listSku, sk *sku.Transacted) error,
) (err error) {
	for list, iterErr := range store.IterAllInventoryLists() {
		if iterErr != nil {
			err = errors.Wrap(iterErr)
			return
		}

		t := list

		if err = f(t, t); err != nil {
			err = errors.Wrapf(
				err,
				"InventoryList: %s",
				t.GetObjectId(),
			)

			return
		}

		iter := store.IterInventoryList(
			t.GetBlobSha(),
		)

		for sk, iterErr := range iter {
			if iterErr != nil {
				err = errors.Wrap(iterErr)
				return
			}

			if err = f(t, sk); err != nil {
				err = errors.Wrapf(
					err,
					"InventoryList: %s",
					t.GetObjectId(),
				)

				return
			}
		}
	}

	return
}
