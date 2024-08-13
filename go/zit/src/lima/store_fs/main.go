package store_fs

import (
	"encoding/gob"
	"strings"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/query_spec"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
)

func init() {
	gob.Register(External{})
}

type Store struct {
	config              sku.Config
	deletedPrinter      interfaces.FuncIter[*fd.FD]
	metadataTextParser  object_metadata.TextParser
	fs_home             fs_home.Home
	fileEncoder         FileEncoder
	ic                  ids.InlineTypeChecker
	fileExtensions      file_extensions.FileExtensions
	dir                 string
	objectFormatOptions object_inventory_format.Options

	dirFDs

	deleteLock sync.Mutex
	deleted    fd.MutableSet
}

func (fs *Store) GetExternalStoreLike() external_store.StoreLike {
	return fs
}

func (fs *Store) DeleteCheckout(col sku.CheckedOutLike) (err error) {
	e := col.GetSkuExternalLike().(*External)

	fs.deleteLock.Lock()
	defer fs.deleteLock.Unlock()

	if err = fs.deleted.Add(&e.FDs.Conflict); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = fs.deleted.Add(e.GetObjectFD()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = fs.deleted.Add(e.GetBlobFD()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.FDs.MutableSetLike.Each(fs.deleted.Add); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (fs *Store) Flush() (err error) {
	deleteOp := DeleteCheckout{}

	if err = deleteOp.Run(
		fs.config.IsDryRun(),
		fs.fs_home,
		fs.deletedPrinter,
		fs.deleted,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	fs.deleted.Reset()

	return
}

func (fs *Store) String() (out string) {
	if iter.Len(
		fs.dirFDs.objects,
		fs.blobs,
	) == 0 {
		return
	}

	sb := &strings.Builder{}
	sb.WriteRune(query_spec.OpGroupOpen)

	hasOne := false

	writeOneIfNecessary := func(v interfaces.Stringer) (err error) {
		if hasOne {
			sb.WriteRune(query_spec.OpOr)
		}

		sb.WriteString(v.String())

		hasOne = true

		return
	}

	fs.dirFDs.objects.Each(
		func(z *FDSet) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.blobs.Each(
		func(z *FDSet) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	sb.WriteRune(query_spec.OpGroupClose)

	out = sb.String()
	return
}

func (s *Store) GetExternalObjectIds() (ks []sku.ExternalObjectId, err error) {
	ks = make([]sku.ExternalObjectId, 0)
	var l sync.Mutex

	if err = s.All(
		func(kfp *FDSet) (err error) {
			l.Lock()
			defer l.Unlock()

			ks = append(ks, kfp)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) GetObjectIdsForDir(
	fd *fd.FD,
) (k []sku.ExternalObjectId, err error) {
	if !fd.IsDir() {
		err = errors.Errorf("not a directory: %q", fd)
		return
	}

	var results []*FDSet

	if results, err = s.dirFDs.processDir(fd.GetPath()); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, r := range results {
		k = append(k, r)
	}

	return
}

// TODO confirm against actual Object Id
func (s *Store) GetObjectIdsForString(v string) (k []sku.ExternalObjectId, err error) {
	if v == "." {
		if k, err = s.GetExternalObjectIds(); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	var fdee *fd.FD

	if fdee, err = fd.MakeFromPath(v, s.fs_home); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fdee.IsDir() {
		if k, err = s.GetObjectIdsForDir(fdee); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		var results []*FDSet

		if _, results, err = s.dirFDs.processFD(fdee); err != nil {
			err = errors.Wrap(err)
			return
		}

		k = make([]sku.ExternalObjectId, 0, len(results))

		for _, r := range results {
			k = append(k, r)
		}
	}

	return
}

func (fs *Store) ContainsSku(m *sku.Transacted) bool {
	return fs.dirFDs.objects.ContainsKey(m.GetObjectId().String())
}

func (fs *Store) GetBlobFDs() fd.Set {
	fds := fd.MakeMutableSet()

	fs.blobs.Each(
		func(fds *FDSet) error {
			return fds.Each(fds.Add)
		},
	)

	return fds
}

func (fs *Store) GetEmptyDirectories() fd.Set {
	fds := fd.MakeMutableSet()
	fs.emptyDirectories.Each(fds.Add)
	return fds
}

func (fs *Store) Get(
	k interfaces.ObjectId,
) (t *FDSet, ok bool) {
	return fs.dirFDs.objects.Get(k.String())
}

func (s *Store) Initialize(esi external_store.Info) (err error) {
	s.externalStoreInfo = esi

	if err = s.dirFDs.processRootDir(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Store) Len() int {
	return iter.Len(
		c.dirFDs.objects,
	)
}
