package store_fs

import (
	"encoding/gob"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
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

type objectsAndBlobs struct {
	unsureZettels interfaces.MutableSetLike[*ObjectIdFDPair]
	objects       interfaces.MutableSetLike[*ObjectIdFDPair]
	blobs         fd.MutableSet
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
	objectsAndBlobs
	emptyDirectories fd.MutableSet

	objectFormatOptions object_inventory_format.Options

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
		fs.objects,
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

	fs.objects.Each(
		func(z *ObjectIdFDPair) (err error) {
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
		func(kfp *ObjectIdFDPair) (err error) {
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

	// TODO implement traversal

	return
}

// TODO confirm against actual Object Id
func (s *Store) GetObjectIdsForString(v string) (k []*ids.ObjectId, err error) {
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
		var oid *ObjectIdFDPair

		if oid, err = s.tryFD(&fd, s.objectsAndBlobs); err != nil {
			err = errors.Wrap(err)
			return
		}

		k = []*ids.ObjectId{&oid.ObjectId}
	}

	return
}

func (fs *Store) ContainsSku(m *sku.Transacted) bool {
	return fs.objects.ContainsKey(m.GetObjectId().String())
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
) (t *ObjectIdFDPair, ok bool) {
	return fs.objects.Get(k.String())
}

func (fs *Store) All(
	f interfaces.FuncIter[*ObjectIdFDPair],
) (err error) {
	wg := iter.MakeErrorWaitGroupParallel()

	iter.ErrorWaitGroupApply(
		wg,
		fs.objects,
		func(e *ObjectIdFDPair) (err error) {
			return f(e)
		},
	)

	iter.ErrorWaitGroupApply(
		wg,
		fs.unsureZettels,
		func(e *ObjectIdFDPair) (err error) {
			return f(e)
		},
	)

	return wg.GetError()
}

func (fs *Store) AllUnsure(
	f interfaces.FuncIter[*ObjectIdFDPair],
) (err error) {
	wg := iter.MakeErrorWaitGroupParallel()

	iter.ErrorWaitGroupApply(
		wg,
		fs.unsureZettels,
		func(e *ObjectIdFDPair) (err error) {
			return f(e)
		},
	)

	return wg.GetError()
}

func (fs *Store) readInputFiles(args ...string) (err error) {
	for _, f := range args {
		f = filepath.Clean(f)

		if filepath.IsAbs(f) {
			if f, err = filepath.Rel(fs.dir, f); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		parts := strings.Split(f, string(filepath.Separator))

		switch len(parts) {
		case 0:

		case 1:
			if err = fs.readNotSecondLevelFile(parts[0]); err != nil {
				err = errors.Wrap(err)
				return
			}

		case 2:
			p := path.Join(parts[len(parts)-2], parts[len(parts)-1])

			if err = fs.readSecondLevelFile(fs.dir, p); err != nil {
				err = errors.Wrap(err)
				return
			}

		default:
			h := path.Join(parts[:len(parts)-3]...)
			p := path.Join(parts[len(parts)-2], parts[len(parts)-1])

			if err = fs.readSecondLevelFile(h, p); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

func (s *Store) Initialize(esi external_store.Info) (err error) {
	s.externalStoreInfo = esi
	return
}

func (s *Store) readAll() (err error) {
	{
		_, err := makeDir(s.dir, s.fileExtensions)
		errors.PanicIfError(err)
	}
	// TODO use walkdir instead
	// check for empty directories
	if err = filepath.WalkDir(
		s.dir,
		func(p string, d fs.DirEntry, in error) (err error) {
			if in != nil {
				err = errors.Wrap(in)
				return
			}

			var rel string

			if rel, err = filepath.Rel(s.dir, p); err != nil {
				err = errors.Wrap(in)
				return
			}

			dir := filepath.Dir(p)
			base := filepath.Base(p)

			if strings.HasPrefix(dir, ".") ||
				strings.HasPrefix(base, ".") ||
				strings.HasPrefix(rel, ".") {
				err = filepath.SkipDir
				return
			}

			if d.IsDir() {
				if strings.HasPrefix(p, ".") {
					err = filepath.SkipDir
				}

				return
			}

			levels := files.DirectoriesRelativeTo(rel)

			if len(levels) == 1 {
				ui.Log().Print("second", rel)
			} else {
				ui.Log().Print("not second", rel)
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var dirs []string

	if dirs, err = files.ReadDirNames(s.dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, d := range dirs {
		if strings.HasPrefix(d, ".") {
			continue
		}

		d2 := path.Join(s.dir, d)

		var fi os.FileInfo

		if fi, err = os.Stat(d); err != nil {
			err = errors.Wrap(err)
			return
		}

		var f *fd.FD

		if f, err = fd.MakeFromFileInfoWithDir(fi, s.dir); err != nil {
			err = errors.Wrap(err)
			return
		}

		if fi.Mode().IsDir() {
			var dirs2 []string

			if dirs2, err = files.ReadDirNames(d2); err != nil {
				err = errors.Wrap(err)
				return
			}

			if len(dirs2) == 0 {
				if err = s.emptyDirectories.Add(f); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			for _, a := range dirs2 {
				if err = s.readSecondLevelFile(d2, a); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

		} else if fi.Mode().IsRegular() {
			if err = s.readNotSecondLevelFile(d); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

func (c *Store) Len() int {
	return iter.Len(
		c.objects,
	)
}

func (s *Store) readNotSecondLevelFile(name string) (err error) {
	if strings.HasPrefix(name, ".") {
		return
	}

	fullPath := path.Join(s.dir, name)

	var fi os.FileInfo

	if fi, err = os.Stat(fullPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !fi.Mode().IsRegular() {
		return
	}

	ext := filepath.Ext(name)
	ext = strings.ToLower(ext)
	ext = strings.TrimSpace(ext)

	var f *fd.FD

	if f, err = fd.MakeFromFileInfoWithDir(fi, s.dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	var oid *ObjectIdFDPair

	if oid, err = s.tryFD(f, s.objectsAndBlobs); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.objects.Add(oid); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
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

func (s *Store) readSecondLevelFile(dir string, name string) (err error) {
	if strings.HasPrefix(name, ".") {
		return
	}

	var fi os.FileInfo

	fullPath := path.Join(dir, name)

	if fi, err = os.Stat(fullPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !fi.Mode().IsRegular() {
		return
	}

	var f *fd.FD

	if f, err = fd.MakeFromFileInfoWithDir(fi, dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	var oid *ObjectIdFDPair

	if oid, err = s.tryFD(f, s.objectsAndBlobs); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.objects.Add(oid); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) tryFD(
	f *fd.FD,
	prior objectsAndBlobs,
) (oidPair *ObjectIdFDPair, err error) {
	depth := f.DepthRelativeTo(s.dir)
	key := f.FileNameSansExt()
	var g genres.Genre
	ext := f.ExtSansDot()
	isConflict := false

	if ext == "conflict" {
		isConflict = true
		ext = fd.ExtSansDot(strings.TrimSuffix(f.GetPath(), f.Ext()))
	}

	switch ext {
	case s.fileExtensions.Zettel:
		g = genres.Zettel

		if depth == 1 {
			key = strings.ToLower(filepath.Join(f.DirBaseOnly(), key))
		} else {
			// recognized
		}

	case s.fileExtensions.Typ:
		g = genres.Type
		key = strings.ToLower(key)

	case s.fileExtensions.Etikett:
		g = genres.Tag
		key = strings.ToLower(key)

	case s.fileExtensions.Kasten:
		g = genres.Repo
		key = strings.ToLower(key)

	default: // blobs
		// TODO
	}

	if isConflict {
	}

	var ok bool

	if prior.objects != nil {
		oidPair, ok = prior.objects.Get(key)
	}

	if !ok {
		oidPair = &ObjectIdFDPair{}

		if err = oidPair.ObjectId.SetWithGenre(key, g); err != nil {
			err = errors.Wrapf(err, "FD: %q", f)
			return
		}
	}

	oidPair.FDs.Object.ResetWith(f)

	return
}
