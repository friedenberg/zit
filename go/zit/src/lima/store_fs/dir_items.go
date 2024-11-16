package store_fs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
)

type fdSetWithError struct {
	error
	*sku.FSItem
}

// TODO support globs and ignores
type dirItems struct {
	root          string
	rootProcessed bool

	interfaces.FileExtensionGetter
	dirLayout             dir_layout.DirLayout
	externalStoreSupplies external_store.Supplies

	probablyCheckedOut      fsItemData
	definitelyNotCheckedOut fsItemData

	errors interfaces.MutableSetLike[fdSetWithError]
}

func makeObjectsWithDir(
	p string,
	fe interfaces.FileExtensionGetter,
	fs_home dir_layout.DirLayout,
) (d dirItems) {
	d.root = p
	d.FileExtensionGetter = fe
	d.dirLayout = fs_home
	d.probablyCheckedOut = makeFSItemData()
	d.definitelyNotCheckedOut = makeFSItemData()
	d.errors = collections_value.MakeMutableValueSet[fdSetWithError](nil)

	return
}

//  __        __    _ _    _
//  \ \      / /_ _| | | _(_)_ __   __ _
//   \ \ /\ / / _` | | |/ / | '_ \ / _` |
//    \ V  V / (_| | |   <| | | | | (_| |
//     \_/\_/ \__,_|_|_|\_\_|_| |_|\__, |
//                                 |___/

func (d *dirItems) walkDir(
	cache map[string]*sku.FSItem,
	dir string,
	pattern string,
) (err error) {
	if err = filepath.WalkDir(
		dir,
		func(p string, de fs.DirEntry, in error) (err error) {
			if in != nil {
				err = errors.Wrap(in)
				return
			}

			if p == d.root {
				return
			}

			if de.Type()&fs.ModeSymlink != 0 {
				if p, err = filepath.EvalSymlinks(p); err != nil {
					err = nil
					return
					// err = errors.Wrap(err)
					// return
				}

				var fi fs.FileInfo

				if fi, err = os.Lstat(p); err != nil {
					err = errors.Wrap(err)
					return
				}

				de = fs.FileInfoToDirEntry(fi)
			}

			if strings.HasPrefix(filepath.Base(p), ".") {
				if de.IsDir() {
					err = filepath.SkipDir
				}

				return
			}

			if pattern != "" {
				var matched bool

				if matched, err = filepath.Match(pattern, p); err != nil {
					err = errors.Wrap(err)
					return
				}

				if !matched {
					return
				}
			}

			if de.IsDir() {
				return
			}

			if _, _, err = d.addPathAndDirEntry(cache, p, de); err != nil {
				err = errors.Wrapf(err, "DirEntry: %s", de)
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

func (d *dirItems) addPathAndDirEntry(
	cache map[string]*sku.FSItem, p string,
	de fs.DirEntry,
) (key string, fds *sku.FSItem, err error) {
	if de.IsDir() {
		return
	}

	var fdee *fd.FD

	if fdee, err = fd.MakeFromPathAndDirEntry(p, de, d.dirLayout); err != nil {
		err = errors.Wrap(err)
		return
	}

	if key, fds, err = d.addFD(cache, fdee); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (d *dirItems) keyForFD(fdee *fd.FD) (key string, err error) {
	var rel string

	if rel, err = filepath.Rel(d.root, fdee.GetPath()); err != nil {
		err = errors.Wrap(err)
		return
	}

	ext := filepath.Ext(rel)
	key = strings.TrimSuffix(rel, ext)

	return
}

func (d *dirItems) addFD(
	cache map[string]*sku.FSItem, f *fd.FD,
) (key string, fds *sku.FSItem, err error) {
	if f.IsDir() {
		return
	}

	if key, err = d.keyForFD(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	if cache == nil {
		fds = &sku.FSItem{
			MutableSetLike: collections_value.MakeMutableValueSet[*fd.FD](nil),
		}

		fds.Add(f)
	} else {
		var ok bool
		fds, ok = cache[key]

		if !ok {
			fds = &sku.FSItem{
				MutableSetLike: collections_value.MakeMutableValueSet[*fd.FD](nil),
			}
		}

		fds.Add(f)
		cache[key] = fds
	}

	return
}

//   ____                              _
//  |  _ \ _ __ ___   ___ ___  ___ ___(_)_ __   __ _
//  | |_) | '__/ _ \ / __/ _ \/ __/ __| | '_ \ / _` |
//  |  __/| | | (_) | (_|  __/\__ \__ \ | | | | (_| |
//  |_|   |_|  \___/ \___\___||___/___/_|_| |_|\__, |
//                                             |___/

func (d *dirItems) processDir(p string) (results []*sku.FSItem, err error) {
	cache := make(map[string]*sku.FSItem)

	results = make([]*sku.FSItem, 0)

	if err = d.walkDir(cache, p, ""); err != nil {
		err = errors.Wrap(err)
		return
	}

	for objectIdString, fds := range cache {
		var someResult []*sku.FSItem

		if someResult, err = d.processFDSet(objectIdString, fds); err != nil {
			err = errors.Wrap(err)
			return
		}

		results = append(results, someResult...)
	}

	return
}

func (d *dirItems) processFD(
	fdee *fd.FD,
) (objectIdString string, fds []*sku.FSItem, err error) {
	cache := make(map[string]*sku.FSItem)

	if objectIdString, err = d.keyForFD(fdee); err != nil {
		err = errors.Wrap(err)
		return
	}

	dir := filepath.Dir(fdee.GetPath())
	pattern := filepath.Join(dir, fmt.Sprintf("%s*", fdee.FileNameSansExt()))

	if err = d.walkDir(cache, dir, pattern); err != nil {
		err = errors.Wrap(err)
		return
	}

	item := cache[objectIdString]

	if item == nil {
		err = errors.Errorf(
			"failed to write FSItem to cache. Cache: %s, Pattern: %s, ObjectId: %s, Dir: %s",
			cache,
			pattern,
			objectIdString,
			dir,
		)

		panic(err)
	}

	if fds, err = d.processFDSet(
		objectIdString,
		item,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (d *dirItems) getFDsForObjectIdString(
	objectIdString string,
) (fds []*sku.FSItem, err error) {
	cache := make(map[string]*sku.FSItem)
	dir := d.dirLayout.Cwd()
	pattern := filepath.Join(dir, fmt.Sprintf("%s*", objectIdString))

	if err = d.walkDir(cache, dir, pattern); err != nil {
		err = errors.Wrap(err)
		return
	}

	item := cache[objectIdString]

	if item == nil {
		err = errors.Errorf(
			"failed to write FSItem to cache. Cache: %s, Pattern: %s, ObjectId: %s, Dir: %s",
			cache,
			pattern,
			objectIdString,
			dir,
		)

		return
	}

	if fds, err = d.processFDSet(
		objectIdString,
		item,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (d *dirItems) processRootDir() (err error) {
	if d.rootProcessed {
		return
	}

	if _, err = d.processDir(d.root); err != nil {
		err = errors.Wrap(err)
		return
	}

	d.rootProcessed = true

	return
}

func (d *dirItems) processFDsOnItem(
	fds *sku.FSItem,
) (blobCount, objectCount int, err error) {
	if err = fds.Each(
		func(f *fd.FD) (err error) {
			ext := f.ExtSansDot()

			switch ext {
			case d.GetFileExtensionZettel():
				fds.ExternalObjectId.SetGenre(genres.Zettel)

			case d.GetFileExtensionType():
				fds.ExternalObjectId.SetGenre(genres.Type)

			case d.GetFileExtensionTag():
				fds.ExternalObjectId.SetGenre(genres.Tag)

			case d.GetFileExtensionRepo():
				fds.ExternalObjectId.SetGenre(genres.Repo)

			case "conflict":
				fds.Conflict.ResetWith(f)
				return

			default: // blobs
				fds.Blob.ResetWith(f)
				blobCount++
				return
			}

			fds.Object.ResetWith(f)
			objectCount++

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (d *dirItems) processFDSet(
	objectIdString string,
	fds *sku.FSItem,
) (results []*sku.FSItem, err error) {
	var recognizedGenre genres.Genre

	{
		recognized := sku.GetTransactedPool().Get()
		defer sku.GetTransactedPool().Put(recognized)

		if err = d.externalStoreSupplies.FuncReadOneInto(
			objectIdString,
			recognized,
		); err != nil {
			if collections.IsErrNotFound(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		} else {
			recognizedGenre = genres.Must(recognized.GetGenre())
		}
	}

	var blobCount, objectCount int

	if blobCount, objectCount, err = d.processFDsOnItem(fds); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fds.ExternalObjectId.GetGenre() != genres.None {
		if blobCount > 1 {
			err = errors.Errorf(
				"several blobs matching object id %q: %q",
				objectIdString,
				fds.MutableSetLike,
			)
		} else if objectCount > 1 {
			err = errors.Errorf(
				"found more than one object: %q",
				fds.MutableSetLike,
			)
		}

		if err != nil {
			if err = d.errors.Add(fdSetWithError{FSItem: fds, error: err}); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	if fds.ExternalObjectId.GetGenre() == genres.None {
		fds.ExternalObjectId.SetGenre(recognizedGenre)
	}

	if fds.ExternalObjectId.GetGenre() == genres.None {
		if results, err = d.addOneOrMoreBlobs(
			fds,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = d.addOneObject(
			objectIdString,
			fds,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		results = []*sku.FSItem{fds}
	}

	return
}

func (d *dirItems) addOneUntracked(
	f *fd.FD,
) (result *sku.FSItem, err error) {
	result = &sku.FSItem{
		MutableSetLike: collections_value.MakeMutableValueSet[*fd.FD](nil),
	}

	result.Blob.ResetWith(f)

	if err = result.MutableSetLike.Add(&result.Blob); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = result.ExternalObjectId.SetBlob(
		d.dirLayout.Rel(f.GetPath()),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = d.definitelyNotCheckedOut.Add(result); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := f.GetSha()

	if sh.IsNull() {
		return
	}

	// TODO try reading as object

	// TODO add sha cache
	key := sh.GetBytes()
	existing, ok := d.definitelyNotCheckedOut.shas[key]

	if !ok {
		existing = collections_value.MakeMutableValueSet[*sku.FSItem](nil)
	}

	if err = existing.Add(result); err != nil {
		err = errors.Wrap(err)
		return
	}

	d.definitelyNotCheckedOut.shas[key] = existing

	return
}

func (d *dirItems) addOneOrMoreBlobs(
	fds *sku.FSItem,
) (results []*sku.FSItem, err error) {
	if fds.MutableSetLike.Len() == 1 {
		var fdsOne *sku.FSItem

		if fdsOne, err = d.addOneUntracked(
			fds.Any(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		fdsOne.ExternalObjectId.SetGenre(genres.Blob)
		results = []*sku.FSItem{fdsOne}

		return
	}

	if err = fds.MutableSetLike.Each(
		func(fd *fd.FD) (err error) {
			var fdsOne *sku.FSItem

			if fdsOne, err = d.addOneUntracked(
				fds.Any(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			results = append(results, fdsOne)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (d *dirItems) addOneObject(
	objectIdString string,
	item *sku.FSItem,
) (err error) {
	g := item.ExternalObjectId.GetGenre()

	if g == genres.Zettel {
		err = item.ExternalObjectId.SetWithGenre(fd.ZettelId(objectIdString), g)
	} else {
		err = item.ExternalObjectId.SetWithGenre(objectIdString, g)
	}

	if err != nil {
		item.ExternalObjectId.SetGenre(item.ExternalObjectId.GetGenre())

		if err = item.ExternalObjectId.SetRaw(objectIdString); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = d.probablyCheckedOut.Add(item); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

//   ___ _                 _   _
//  |_ _| |_ ___ _ __ __ _| |_(_) ___  _ __
//   | || __/ _ \ '__/ _` | __| |/ _ \| '_ \
//   | || ||  __/ | | (_| | |_| | (_) | | | |
//  |___|\__\___|_|  \__,_|\__|_|\___/|_| |_|
//

func (d *dirItems) All(
	f interfaces.FuncIter[*sku.FSItem],
) (err error) {
	wg := quiter.MakeErrorWaitGroupParallel()

	quiter.ErrorWaitGroupApply(wg, d.probablyCheckedOut, f)
	quiter.ErrorWaitGroupApply(wg, d.definitelyNotCheckedOut, f)

	return wg.GetError()
}
