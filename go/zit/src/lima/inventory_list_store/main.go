package inventory_list_store

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
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
	"code.linenisgreat.com/zit/go/zit/src/india/inventory_list_fmt"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt_debug"
)

type Store struct {
	dirLayout dir_layout.DirLayout
	ls        interfaces.LockSmith
	sv        interfaces.StoreVersion
	of        interfaces.ObjectIOFactory
	af        interfaces.BlobIOFactory
	clock     ids.Clock
	pool      interfaces.Pool[sku.List, *sku.List]

	object_format object_inventory_format.Format
	options       object_inventory_format.Options
	box           *box_format.Box

	sku.ListFormat
}

func (s *Store) Initialize(
	dirLayout dir_layout.DirLayout,
	ls interfaces.LockSmith,
	sv interfaces.StoreVersion,
	of interfaces.ObjectIOFactory,
	af interfaces.BlobIOFactory,
	pmf object_inventory_format.Format,
	clock ids.Clock,
	box *box_format.Box,
) (err error) {
	p := pool.MakePool(nil, func(a *sku.List) { sku.ResetterList.Reset(a) })

	op := object_inventory_format.Options{Tai: true}

	*s = Store{
		dirLayout:     dirLayout,
		ls:            ls,
		sv:            sv,
		of:            of,
		af:            af,
		pool:          p,
		clock:         clock,
		object_format: pmf,
		options:       op,
		box:           box,
	}

	s.ListFormat = s.FormatForVersion(sv)

	return
}

func (s *Store) Flush() (err error) {
	wg := quiter.MakeErrorWaitGroupParallel()
	return wg.GetError()
}

func (s *Store) FormatForVersion(sv interfaces.StoreVersion) sku.ListFormat {
	v := sv.GetInt()

	switch {
	case v <= 6:
		return inventory_list_fmt.FormatOld{
			Factory: inventory_list_fmt.Factory{
				Format:  s.object_format,
				Options: s.options,
			},
		}

	default:
		return inventory_list_fmt.FormatNew{
			Box: s.box,
		}
	}
}

func (s *Store) Create(
	o *sku.List,
	bez descriptions.Description,
) (t *sku.Transacted, err error) {
	if o.Len() == 0 {
		return
	}

	if !s.ls.IsAcquired() {
		err = file_lock.ErrLockRequired{
			Operation: "create bestandsaufnahme",
		}

		return
	}

	t = sku.GetTransactedPool().Get()

	if o.Len() > 0 {
		var wc interfaces.ShaWriteCloser

		if wc, err = s.dirLayout.BlobWriter(); err != nil {
			err = errors.Wrap(err)
			return
		}

		func() {
			defer errors.DeferredCloser(&err, wc)

			if _, err = s.WriteInventoryListBlob(o, wc); err != nil {
				err = errors.Wrap(err)
				return
			}
		}()

		sh := wc.GetShaLike()

		if _, err = s.GetBlob(sh); err != nil {
			err = errors.Wrapf(err, "Blob Sha: %q", sh)
			return
		}

		t.SetBlobSha(sh)
	}

	t.Metadata.Description = bez
	tai := s.clock.GetTai()

	if err = t.ObjectId.SetWithIdLike(tai); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.SetTai(tai)

	var wc interfaces.ShaWriteCloser

	if wc, err = s.of.ObjectWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	func() {
		defer errors.DeferredCloser(&err, wc)

		if _, err = s.ListFormat.WriteInventoryListObject(
			t,
			wc,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}()

	sh := wc.GetShaLike()

	ui.Log().Printf(
		"saving Bestandsaufnahme with tai: %s -> %s",
		t.GetObjectId().GetGenre(),
		sh,
	)

	if err = t.CalculateObjectShas(); err != nil {
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

	if _, o, err = s.ReadInventoryListObject(or); err != nil {
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

	if err = s.StreamInventoryListBlobSkus(
		rc,
		f,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) readInventoryListBlob(
	blobSha interfaces.Sha,
	a *sku.List,
) (err error) {
	var rc interfaces.ShaReadCloser

	if rc, err = s.af.BlobReader(blobSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, rc)

	if err = s.StreamInventoryListBlobSkus(
		rc,
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

func (s *Store) GetBlob(blobSha interfaces.Sha) (a *sku.List, err error) {
	a = sku.MakeList()

	var rc interfaces.ShaReadCloser

	if rc, err = s.af.BlobReader(blobSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, rc)

	if err = inventory_list_fmt.ReadInventoryListBlob(
		s.ListFormat,
		rc,
		a,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO validate sha

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

	if p, err = s.dirLayout.DirObjectGenre(
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
