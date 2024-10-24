package inventory_list

import (
	"fmt"
	"sort"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_lock"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt_debug"
)

type Format = blob_store.Format[
	InventoryList,
	*InventoryList,
]

type Store struct {
	fs_home       fs_home.Home
	ls            interfaces.LockSmith
	sv            interfaces.StoreVersion
	of            interfaces.ObjectIOFactory
	af            interfaces.BlobIOFactory
	clock         ids.Clock
	pool          interfaces.Pool[InventoryList, *InventoryList]
	object_format object_inventory_format.Format
	options       object_inventory_format.Options
	format

	versionedFormat
}

func (s *Store) Initialize(
	fs_home fs_home.Home,
	ls interfaces.LockSmith,
	sv interfaces.StoreVersion,
	of interfaces.ObjectIOFactory,
	af interfaces.BlobIOFactory,
	pmf object_inventory_format.Format,
	clock ids.Clock,
	box *box_format.Box,
) (err error) {
	p := pool.MakePool(nil, func(a *InventoryList) { Resetter.Reset(a) })

	op := object_inventory_format.Options{Tai: true}
	fa := MakeFormat(sv, op)

	*s = Store{
		fs_home:       fs_home,
		ls:            ls,
		sv:            sv,
		of:            of,
		af:            af,
		pool:          p,
		clock:         clock,
		object_format: pmf,
		options:       op,
		format:        fa,
	}

	v := sv.GetInt()

	switch {
	case v <= 6 || true:
		s.versionedFormat = versionedFormatOld{
			object_format: pmf,
			options:       op,
			format:        fa,
		}

	default:
		s.versionedFormat = versionedFormatNew{
			box: box,
		}
	}

	return
}

func (s *Store) Flush() (err error) {
	wg := quiter.MakeErrorWaitGroupParallel()
	return wg.GetError()
}

func (s *Store) Create(
	o *InventoryList,
	bez descriptions.Description,
) (t *sku.Transacted, err error) {
	if !s.ls.IsAcquired() {
		err = file_lock.ErrLockRequired{
			Operation: "create bestandsaufnahme",
		}

		return
	}

	if o.Len() == 0 {
		err = errors.Wrap(ErrEmpty)
		return
	}

	var sh *sha.Sha

	if sh, err = s.writeInventoryListBlob(o, s.fs_home.BlobWriter); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = s.GetBlob(sh); err != nil {
		err = errors.Wrapf(err, "Blob Sha: %q", sh)
		return
	}

	t = sku.GetTransactedPool().Get()

	sku.TransactedResetter.Reset(t)
	t.Metadata.Description = bez
	t.SetBlobSha(sh)
	tai := s.clock.GetTai()

	if err = t.ObjectId.SetWithIdLike(tai); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.SetTai(tai)

	if sh, err = s.versionedFormat.writeInventoryListObject(
		t,
		s.of.ObjectWriter,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	ui.Log().Printf(
		"saving Bestandsaufnahme with tai: %s -> %s",
		t.GetObjectId().GetGenre(),
		sh,
	)

	t.SetObjectSha(sh)

	return
}

func (s *Store) readOnePath(p string) (o *sku.Transacted, err error) {
	var sh *sha.Sha

	if sh, err = sha.MakeShaFromPath(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if o, err = s.ReadOne(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = o.CalculateObjectShas(); err != nil {
		if errors.Is(err, object_inventory_format.ErrEmptyTai) {
			var t ids.Tai
			err1 := t.Set(o.ObjectId.String())

			if err1 != nil {
				err = errors.Wrapf(err, "%s", sku_fmt_debug.StringTaiGenreObjectIdShaBlob(o))
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

func (s *Store) ReadOne(
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

	if _, o, err = s.readInventoryListObject(or); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) StreamInventoryList(
	blobSha interfaces.Sha,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	if err = s.streamInventoryListBlobSkus(
		s.af.BlobReader,
		blobSha,
		f,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) readInventoryListBlob(
	rf func(interfaces.ShaGetter) (interfaces.ShaReadCloser, error),
	blobSha interfaces.Sha,
	a *InventoryList,
) (err error) {
	if err = s.streamInventoryListBlobSkus(
		rf,
		blobSha,
		func(sk *sku.Transacted) (err error) {
			if err = a.Add(sk); err != nil {
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

func (s *Store) GetBlob(blobSha interfaces.Sha) (a *InventoryList, err error) {
	a = MakeInventoryList()
	err = s.readInventoryListBlob(s.af.BlobReader, blobSha, a)
	return
}

func (s *Store) ReadLast() (max *sku.Transacted, err error) {
	l := &sync.Mutex{}

	var maxSku sku.Transacted

	if err = s.ReadAll(
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
				"did not find last Bestandsaufnahme: %#v",
				max.GetMetadata(),
			),
		)
	}

	return
}

func (s *Store) ReadAll(
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	var p string

	if p, err = s.fs_home.DirObjectGenre(
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

	if err = s.ReadAll(
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
	f func(besty, sk *sku.Transacted) error,
) (err error) {
	if err = s.ReadAll(
		func(t *sku.Transacted) (err error) {
			if err = s.StreamInventoryList(
				t.GetBlobSha(),
				func(sk *sku.Transacted) (err error) {
					return f(t, sk)
				},
			); err != nil {
				err = errors.Wrapf(
					err,
					"Bestandsaufnahme: %s",
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
