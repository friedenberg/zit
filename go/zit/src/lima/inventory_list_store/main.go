package inventory_list_store

import (
	"fmt"
	"iter"
	"sort"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_lock"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
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

	envRepo        env_repo.Env
	ls             interfaces.LockSmith
	sv             interfaces.StoreVersion
	of             interfaces.ObjectIOFactory
	af             interfaces.BlobStore
	clock          ids.Clock
	typedBlobStore typed_blob_store.InventoryList

	object_format object_inventory_format.Format
	options       object_inventory_format.Options
	box           *box_format.BoxTransacted

	blobType ids.Type
}

func (s *Store) Initialize(
	envRepo env_repo.Env,
	clock ids.Clock,
	typedBlobStore typed_blob_store.InventoryList,
) (err error) {
	op := object_inventory_format.Options{Tai: true}

	*s = Store{
		envRepo: envRepo,
		ls:      envRepo.GetLockSmith(),
		sv:      envRepo.GetStoreVersion(),
		of:      envRepo.ObjectReaderWriterFactory(genres.InventoryList),
		af:      envRepo,
		clock:   clock,
		box: box_format.MakeBoxTransactedArchive(
			envRepo,
			options_print.V0{}.WithPrintTai(true),
		),
		options:        op,
		typedBlobStore: typedBlobStore,
	}

	v := s.sv.GetInt()

	switch {
	case v <= 6:
		s.blobType = ids.MustType(builtin_types.InventoryListTypeV0)

	default:
		s.blobType = ids.MustType(builtin_types.InventoryListTypeV1)
	}

	return
}

func (s *Store) GetEnv() env_ui.Env {
	return s.GetRepoLayout()
}

func (s *Store) GetImmutableConfig() config_immutable_io.ConfigLoaded {
	return s.GetRepoLayout().GetConfig()
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

func (s *Store) GetRepoLayout() env_repo.Env {
	return s.envRepo
}

func (s *Store) GetBlobStore() interfaces.BlobStore {
	return &s.envRepo
}

func (s *Store) GetInventoryListStore() sku.InventoryListStore {
	return s
}

func (s *Store) Create(
	skus *sku.List,
	description descriptions.Description,
) (t *sku.Transacted, err error) {
	if skus.Len() == 0 {
		return
	}

	if !s.ls.IsAcquired() {
		err = file_lock.ErrLockRequired{
			Operation: "create inventory list",
		}

		return
	}

	t = sku.GetTransactedPool().Get()

	t.Metadata.Type = s.blobType
	t.Metadata.Description = description

	tai := s.GetTai()

	if err = t.ObjectId.SetWithIdLike(tai); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.SetTai(tai)

	if err = s.WriteInventoryListBlob(t, skus); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.WriteInventoryListObject(t); err != nil {
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
func (s *Store) WriteInventoryListObject(t *sku.Transacted) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	var wc interfaces.ShaWriteCloser

	// TODO also write to inventory_list_log

	if wc, err = s.of.ObjectWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, wc)

	t.Metadata.Type = s.blobType

	if _, err = s.typedBlobStore.WriteObjectToWriter(
		t,
		wc,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = t.CalculateObjectShas(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ui.Log().Printf(
		"saved inventory list: %q",
		sku.String(t),
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
		s.FormatForVersion(s.sv),
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
			s.GetRepoLayout().GetEnv(),
			s.GetRepoLayout(),
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

	if or, err = s.of.ObjectReader(&sh); err != nil {
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

func (s *Store) StreamInventoryList(
	blobSha interfaces.Sha,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	var rc interfaces.ShaReadCloser

	if rc, err = s.af.BlobReader(blobSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, rc)

	if err = s.typedBlobStore.StreamInventoryListBlobSkusFromReader(
		s.blobType,
		rc,
		f,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadLast() (max *sku.Transacted, err error) {
	var maxSku sku.Transacted

	for listOrError := range s.AllInventoryLists() {
		if listOrError.Error != nil {
			err = errors.Wrap(listOrError.Error)
			return
		}

		b := listOrError.Element

		if sku.TransactedLessor.LessPtr(&maxSku, b) {
			sku.TransactedResetter.ResetWith(&maxSku, b)
		}
	}

	max = &maxSku

	if max.GetObjectSha().IsNull() {
		panic(
			fmt.Sprintf(
				"did not find last inventory list: %#v",
				max.GetMetadata(),
			),
		)
	}

	return
}

func (s *Store) AllInventoryLists() iter.Seq[quiter.ElementOrError[*sku.Transacted]] {
	var p string

	{
		var err error

		if p, err = s.envRepo.DirObjectGenre(
			genres.InventoryList,
		); err != nil {
			err = errors.Wrap(err)
			return func(yield func(quiter.ElementOrError[*sku.Transacted]) bool) {
				yield(quiter.ElementOrError[*sku.Transacted]{Error: err})
			}
		}
	}

	return func(yield func(quiter.ElementOrError[*sku.Transacted]) bool) {
		for pathOrError := range files.DirNamesLevel2(p) {
			if pathOrError.Error != nil {
				if !yield(
					quiter.ElementOrError[*sku.Transacted]{Error: pathOrError.Error},
				) {
					return
				}
			}

			var decodedList *sku.Transacted

			{
				var err error

				if decodedList, err = s.readOnePath(
					pathOrError.Element,
				); err != nil {
					if !yield(
						quiter.ElementOrError[*sku.Transacted]{Error: errors.Wrap(err)},
					) {
						return
					}
				}
			}

			if !yield(quiter.ElementOrError[*sku.Transacted]{Element: decodedList}) {
				return
			}
		}
	}
}

func (s *Store) ReadAllSorted(
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	var skus []*sku.Transacted

	for listOrError := range s.AllInventoryLists() {
		if listOrError.Error != nil {
			err = errors.Wrap(listOrError.Error)
			return
		}

		skus = append(skus, listOrError.Element)
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

func (s *Store) ReadAllSkus(
	f func(listSku, sk *sku.Transacted) error,
) (err error) {
	for listOrError := range s.AllInventoryLists() {
		if listOrError.Error != nil {
			err = errors.Wrap(listOrError.Error)
			return
		}

		t := listOrError.Element

		if err = f(t, t); err != nil {
			err = errors.Wrapf(
				err,
				"InventoryList: %s",
				t.GetObjectId(),
			)

			return
		}

		if err = s.StreamInventoryList(
			t.GetBlobSha(),
			func(sk *sku.Transacted) (err error) {
				return f(t, sk)
			},
		); err != nil {
			err = errors.Wrapf(
				err,
				"InventoryList: %s",
				t.GetObjectId(),
			)

			return
		}
	}

	return
}
