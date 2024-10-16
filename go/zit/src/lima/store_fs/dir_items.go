package store_fs

import (
	"io/fs"
	"path/filepath"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
)

type fdSetWithError struct {
	error
	*Item
}

// TODO support globs and ignores
type dirItems struct {
	root          string
	rootProcessed bool

	file_extensions.FileExtensions
	fs_home               fs_home.Home
	externalStoreSupplies external_store.Supplies

	objects         interfaces.MutableSetLike[*Item]
	blobs           interfaces.MutableSetLike[*Item]
	shasToBlobFDs   map[sha.Bytes]interfaces.MutableSetLike[*Item]
	shasToObjectFDs map[sha.Bytes]interfaces.MutableSetLike[*Item]

	errors interfaces.MutableSetLike[fdSetWithError]
}

func makeObjectsWithDir(
	p string,
	fe file_extensions.FileExtensions,
	fs_home fs_home.Home,
) (d dirItems) {
	d.root = p
	d.FileExtensions = fe
	d.fs_home = fs_home
	d.objects = collections_value.MakeMutableValueSet[*Item](nil)
	d.blobs = collections_value.MakeMutableValueSet[*Item](nil)
	d.shasToBlobFDs = make(map[sha.Bytes]interfaces.MutableSetLike[*Item])
	d.shasToObjectFDs = make(map[sha.Bytes]interfaces.MutableSetLike[*Item])
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
	cache map[string]*Item,
	p string,
) (err error) {
	if err = filepath.WalkDir(
		p,
		func(p string, de fs.DirEntry, in error) (err error) {
			if in != nil {
				err = errors.Wrap(in)
				return
			}

			if p == d.root {
				return
			}

			if strings.HasPrefix(filepath.Base(p), ".") {
				if de.IsDir() {
					err = filepath.SkipDir
				}

				return
			}

			// if de.IsDir() {
			// 	return
			// }

			if _, _, err = d.addPathAndDirEntry(cache, p, de); err != nil {
				err = errors.Wrap(in)
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
	cache map[string]*Item,
	p string,
	de fs.DirEntry,
) (key string, fds *Item, err error) {
	if de.IsDir() {
		return
	}

	var fdee *fd.FD

	if fdee, err = fd.MakeFromPathAndDirEntry(p, de, d.fs_home); err != nil {
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
	cache map[string]*Item,
	f *fd.FD,
) (key string, fds *Item, err error) {
	if f.IsDir() {
		return
	}

	if key, err = d.keyForFD(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	if cache == nil {
		fds = &Item{
			MutableSetLike: collections_value.MakeMutableValueSet[*fd.FD](nil),
		}

		fds.Add(f)
	} else {
		var ok bool
		fds, ok = cache[key]

		if !ok {
			fds = &Item{
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

func (d *dirItems) processDir(p string) (results []*Item, err error) {
	cache := make(map[string]*Item)

	results = make([]*Item, 0)

	if err = d.walkDir(cache, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	for objectIdString, fds := range cache {
		var someResult []*Item

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
) (objectIdString string, fds []*Item, err error) {
	cache := make(map[string]*Item)

	if objectIdString, err = d.keyForFD(fdee); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := filepath.Dir(fdee.GetPath())

	// TODO add filter for just matching fdee
	if err = d.walkDir(cache, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fds, err = d.processFDSet(objectIdString, cache[objectIdString]); err != nil {
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
	fds *Item,
) (blobCount, objectCount int, err error) {
	if err = fds.Each(
		func(f *fd.FD) (err error) {
			ext := f.ExtSansDot()

			switch ext {
			case d.Zettel:
				fds.SetGenre(genres.Zettel)

			case d.Typ:
				fds.SetGenre(genres.Type)

			case d.Etikett:
				fds.SetGenre(genres.Tag)

			case d.Kasten:
				fds.SetGenre(genres.Repo)

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
	fds *Item,
) (results []*Item, err error) {
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

	if fds.GetGenre() != genres.Unknown {
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
			if err = d.errors.Add(fdSetWithError{Item: fds, error: err}); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	if fds.GetGenre() == genres.Unknown {
		fds.ObjectId.SetGenre(recognizedGenre)
	}

	if fds.GetGenre() == genres.Unknown {
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

		results = []*Item{fds}
	}

	return
}

func (d *dirItems) addOneBlob(
	f *fd.FD,
) (result *Item, err error) {
	result = &Item{
		MutableSetLike: collections_value.MakeMutableValueSet[*fd.FD](nil),
	}

	result.Blob.ResetWith(f)

	if err = result.MutableSetLike.Add(&result.Blob); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = result.ObjectId.SetLeft(
		f.GetPath(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = d.blobs.Add(result); err != nil {
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
	existing, ok := d.shasToBlobFDs[key]

	if !ok {
		existing = collections_value.MakeMutableValueSet[*Item](nil)
	}

	if err = existing.Add(result); err != nil {
		err = errors.Wrap(err)
		return
	}

	d.shasToBlobFDs[key] = existing

	return
}

func (d *dirItems) addOneOrMoreBlobs(
	fds *Item,
) (results []*Item, err error) {
	if fds.MutableSetLike.Len() == 1 {
		var fdsOne *Item

		if fdsOne, err = d.addOneBlob(
			fds.Any(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		results = []*Item{fdsOne}

		return
	}

	if err = fds.MutableSetLike.Each(
		func(fd *fd.FD) (err error) {
			var fdsOne *Item

			if fdsOne, err = d.addOneBlob(
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
	fds *Item,
) (err error) {
	g := fds.GetGenre()
	if g == genres.Zettel {
		err = fds.ObjectId.SetWithGenre(fd.ZettelId(objectIdString), g)
	} else {
		err = fds.ObjectId.SetWithGenre(objectIdString, g)
	}

	if err != nil {
		fds.ObjectId.SetGenre(fds.GetGenre())

		if err = fds.ObjectId.SetRaw(objectIdString); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = d.objects.Add(fds); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (d *dirItems) ConsolidateDuplicateBlobs() (err error) {
	replacement := collections_value.MakeMutableValueSet[*Item](nil)

	for _, fds := range d.shasToBlobFDs {
		if fds.Len() == 1 {
			replacement.Add(fds.Any())
		}

		sorted := quiter.ElementsSorted(
			fds,
			func(a, b *Item) bool {
				return a.ObjectId.String() < b.ObjectId.String()
			},
		)

		top := sorted[0]

		for _, other := range sorted[1:] {
			other.MutableSetLike.Each(top.MutableSetLike.Add)
		}

		replacement.Add(top)
	}

	// TODO make less leaky
	d.blobs = replacement

	return
}

//   ___ _                 _   _
//  |_ _| |_ ___ _ __ __ _| |_(_) ___  _ __
//   | || __/ _ \ '__/ _` | __| |/ _ \| '_ \
//   | || ||  __/ | | (_| | |_| | (_) | | | |
//  |___|\__\___|_|  \__,_|\__|_|\___/|_| |_|
//

func (d *dirItems) OnlyObjects(
	f interfaces.FuncIter[*Item],
) (err error) {
	wg := quiter.MakeErrorWaitGroupParallel()

	quiter.ErrorWaitGroupApply(wg, d.objects, f)

	return wg.GetError()
}

func (d *dirItems) All(
	f interfaces.FuncIter[*Item],
) (err error) {
	wg := quiter.MakeErrorWaitGroupParallel()

	quiter.ErrorWaitGroupApply(wg, d.objects, f)
	quiter.ErrorWaitGroupApply(wg, d.blobs, f)

	return wg.GetError()
}
