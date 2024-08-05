package store_fs

import (
	"encoding/gob"
	"path"
	"strings"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
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

// TODO support globs and ignores
type Store struct {
	config             sku.Config
	deletedPrinter     interfaces.FuncIter[*fd.FD]
	externalStoreInfo  external_store.Info
	metadataTextParser object_metadata.TextParser
	fs_home            fs_home.Home
	fileEncoder        FileEncoder
	ic                 ids.InlineTypeChecker
	fileExtensions     file_extensions.FileExtensions
	dir                string
	emptyDirectories   fd.MutableSet

	objectFormatOptions object_inventory_format.Options

	objects

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

// must accept directories
func (fs *Store) MarkUnsureBlob(f *fd.FD) (err error) {
	if f.IsDir() {
		// TODO handle recursion
		return
	}

	if f, err = fd.MakeFromFileFromFD(f, fs.fs_home); err != nil {
		err = errors.Wrapf(err, "%q", f)
		return
	}

	if err = fs.blobs.Add(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (fs *Store) String() (out string) {
	if iter.Len(
		fs.objects.objects,
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

	fs.objects.objects.Each(
		func(z *FDSet) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.blobs.Each(
		func(z *fd.FD) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	sb.WriteRune(query_spec.OpGroupClose)

	out = sb.String()
	return
}

func (s *Store) GetExternalObjectIds() (ks interfaces.SetLike[*ids.ObjectId], err error) {
	ksm := collections_value.MakeMutableValueSet[*ids.ObjectId](nil)
	ks = ksm
	var l sync.Mutex

	if err = s.All(
		func(kfp *FDSet) (err error) {
			kc := kfp.ObjectId.Clone()

			l.Lock()
			defer l.Unlock()

			if err = ksm.Add(kc); err != nil {
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

func (s *Store) GetObjectIdsForDir(fd *fd.FD) (k []*ids.ObjectId, err error) {
	if !fd.IsDir() {
		err = errors.Errorf("not a directory: %q", fd)
		return
	}

	if err = s.objects.walkDir(
		fd.GetPath(),
		func(fds *FDSet) (err error) {
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.objects.ReadObjectsAndBlobs(
		func(fds *FDSet) (err error) {
			k = append(k, fds.ObjectId.Clone())
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO confirm against actual Object Id
func (s *Store) GetObjectIdsForString(v string) (k []*ids.ObjectId, err error) {
	if v == "." {
		v = s.dir
	}

	var fd fd.FD

	if err = fd.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fd.IsDir() {
		if k, err = s.GetObjectIdsForDir(&fd); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		var objectIdString string
		var fds *FDSet

		if objectIdString, fds, err = s.addFD(&fd); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = s.processFDSet(objectIdString, fds); err != nil {
			err = errors.Wrap(err)
			return
		}

		k = []*ids.ObjectId{&fds.ObjectId}
	}

	return
}

func (fs *Store) ContainsSku(m *sku.Transacted) bool {
	return fs.objects.objects.ContainsKey(m.GetObjectId().String())
}

func (fs *Store) GetBlobFDs() fd.Set {
	fds := fd.MakeMutableSet()

	fs.blobs.Each(fds.Add)

	return fds
}

func (fs *Store) GetUnsureBlobs() fd.Set {
	fds := fd.MakeMutableSet()
	fs.blobs.Each(fds.Add)
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
	return fs.objects.objects.Get(k.String())
}

func (fs *Store) All(
	f interfaces.FuncIter[*FDSet],
) (err error) {
	wg := iter.MakeErrorWaitGroupParallel()

	iter.ErrorWaitGroupApply(wg, fs.objects.objects, f)

	iter.ErrorWaitGroupApply(wg, fs.unsureZettels, f)

	return wg.GetError()
}

func (fs *Store) AllUnsure(
	f interfaces.FuncIter[*FDSet],
) (err error) {
	wg := iter.MakeErrorWaitGroupParallel()

	iter.ErrorWaitGroupApply(
		wg,
		fs.unsureZettels,
		func(e *FDSet) (err error) {
			return f(e)
		},
	)

	return wg.GetError()
}

func (s *Store) Initialize(esi external_store.Info) (err error) {
	s.externalStoreInfo = esi
	return
}

func (s *Store) readAll() (err error) {
	if err = s.objects.walkRootDir(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.objects.ReadObjectsAndBlobs(
		func(oidPair *FDSet) (err error) {
			if err = s.objects.objects.Add(oidPair); err != nil {
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

func (c *Store) Len() int {
	return iter.Len(
		c.objects.objects,
	)
}

func (fs *Store) addUnsureBlob(dir, name string) (err error) {
	var ut *fd.FD

	fullPath := name

	if dir != "" {
		fullPath = path.Join(dir, fullPath)
	}

	if ut, err = fd.MakeFromPathWithBlobWriterFactory(
		fullPath,
		fs.fs_home,
	); err != nil {
		err = errors.Wrapf(err, "Dir: %q, Name: %q", dir, name)
		return
	}

	err = fs.blobs.Add(ut)

	return
}
