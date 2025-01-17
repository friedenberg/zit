package inventory_list_store

import (
	"fmt"
	"sort"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_lock"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	"code.linenisgreat.com/zit/go/zit/src/kilo/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/lima/blob_store"
)

type Store struct {
	lock sync.Mutex

	repoLayout repo_layout.Layout
	ls         interfaces.LockSmith
	sv         interfaces.StoreVersion
	of         interfaces.ObjectIOFactory
	af         interfaces.BlobStore
	clock      ids.Clock
	blobStore  blob_store.InventoryList

	object_format object_inventory_format.Format
	options       object_inventory_format.Options
	box           *box_format.BoxTransacted

	blobType ids.Type
}

func (s *Store) Initialize(
	repoLayout repo_layout.Layout,
	pmf object_inventory_format.Format,
	clock ids.Clock,
	blobStore blob_store.InventoryList,
) (err error) {
	op := object_inventory_format.Options{Tai: true}

	*s = Store{
		repoLayout:    repoLayout,
		ls:            repoLayout.GetLockSmith(),
		sv:            repoLayout.GetStoreVersion(),
		of:            repoLayout.ObjectReaderWriterFactory(genres.InventoryList),
		af:            repoLayout,
		clock:         clock,
		object_format: pmf,
		box: box_format.MakeBoxTransactedArchive(
			repoLayout.Env,
			options_print.V0{}.WithPrintTai(true),
		),
		options:   op,
		blobStore: blobStore,
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

func (s *Store) GetEnv() *env.Env {
	return s.GetRepoLayout().Env
}

func (s *Store) GetRepoType() repo_type.Type {
	return s.GetRepoLayout().GetConfig().GetRepoType()
}

func (u *Store) GetStoreVersion() interfaces.StoreVersion {
	return u.GetRepoLayout().GetConfig().GetStoreVersion()
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

func (s *Store) GetRepoLayout() repo_layout.Layout {
	return s.repoLayout
}

func (s *Store) GetBlobStore() interfaces.BlobStore {
	return &s.repoLayout
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

	if wc, err = s.repoLayout.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, wc)

	if _, err = s.blobStore.WriteBlobToWriter(
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

	if _, err = s.blobStore.WriteObjectToWriter(
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

		if _, err = repo_layout.CopyBlobIfNecessary(
			s.GetRepoLayout().GetEnv(),
			s.GetRepoLayout(),
			bs,
			sk.GetBlobSha(),
		); err != nil {
			if errors.Is(err, &dir_layout.ErrAlreadyExists{}) {
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
		err = errors.Wrap(err)
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

	if o, err = s.blobStore.ReadInventoryListObject(
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

	if err = s.blobStore.StreamInventoryListBlobSkusFromReader(
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
	l := &sync.Mutex{}

	var maxSku sku.Transacted

	if err = s.ReadAllInventoryLists(
		func(b *sku.Transacted) (err error) {
			l.Lock()
			defer l.Unlock()

			if sku.TransactedLessor.LessPtr(&maxSku, b) {
				sku.TransactedResetter.ResetWith(&maxSku, b)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
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

func (s *Store) ReadAllInventoryLists(
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	var p string

	if p, err = s.repoLayout.DirObjectGenre(
		genres.InventoryList,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = files.ReadDirNamesLevel2(
		func(p string) (err error) {
			var o *sku.Transacted

			if o, err = s.readOnePath(p); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = f(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
		p,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadAllSorted(
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	var skus []*sku.Transacted

	if err = s.ReadAllInventoryLists(
		func(o *sku.Transacted) (err error) {
			skus = append(skus, o)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
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
	if err = s.ReadAllInventoryLists(
		func(t *sku.Transacted) (err error) {
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

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
