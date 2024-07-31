package store_fs

import (
	"io/fs"
	"path/filepath"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
)

type objects struct {
	root string
	file_extensions.FileExtensions
	fds map[string]*FDSet

	unsureZettels interfaces.MutableSetLike[*FDSet]
	objects       interfaces.MutableSetLike[*FDSet]
	blobs         fd.MutableSet
}

func makeObjectsWithDir(
	p string,
	fe file_extensions.FileExtensions,
) (d objects) {
	d.root = p
	d.FileExtensions = fe
	d.fds = make(map[string]*FDSet)
	d.objects = collections_value.MakeMutableValueSet[*FDSet](nil)
	d.unsureZettels = collections_value.MakeMutableValueSet[*FDSet](nil)
	d.blobs = collections_value.MakeMutableValueSet[*fd.FD](nil)

	return
}

func (d *objects) walkRootDir() (err error) {
	if err = d.walkDir(d.root, nil); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (d *objects) walkDir(
	p string,
	f interfaces.FuncIter[*FDSet],
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

			var fdee *fd.FD

			if fdee, err = fd.MakeFromPathAndDirEntry(p, de); err != nil {
				err = errors.Wrap(err)
				return
			}

			var fds *FDSet

			if _, fds, err = d.addFD(fdee); err != nil {
				err = errors.Wrap(err)
				return
			}

			if f == nil || fds == nil {
				return
			}

			if err = f(fds); err != nil {
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

func (d *objects) addFD(
	f *fd.FD,
) (key string, fds *FDSet, err error) {
	if f.IsDir() {
		return
	}

	var rel string

	if rel, err = filepath.Rel(d.root, f.GetPath()); err != nil {
		err = errors.Wrap(err)
		return
	}

	ext := filepath.Ext(rel)
	key = strings.TrimSuffix(rel, ext)

	var ok bool
	fds, ok = d.fds[key]

	if !ok {
		fds = &FDSet{
			MutableSetLike: collections_value.MakeMutableValueSet[*fd.FD](nil),
		}

		d.fds[key] = fds
	}

	fds.Add(f)
	d.fds[key] = fds

	return
}

func (d *objects) processFDSet(objectIdString string, fds *FDSet) (err error) {
	var blobCount, objectCount int

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
				d.blobs.Add(f)
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

	if fds.GetGenre() != genres.Unknown {
		if blobCount > 1 {
			err = errors.Errorf(
				"several blobs matching object id %q: %q",
				objectIdString,
				fds.MutableSetLike,
			)

			return
		} else if objectCount > 1 {
			err = errors.Errorf(
				"found more than one object: %q",
				fds.MutableSetLike,
			)

			return
		}
	}

	if fds.GetGenre() == genres.Unknown {
		ui.Log().Print(fds.GetGenre())
		if err = fds.ObjectId.SetWithGenre(
			objectIdString,
			genres.Blob,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
		ui.Log().Print(fds.GetGenre())
	} else {
		if err = fds.ObjectId.SetWithGenre(
			objectIdString,
			fds.GetGenre(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	ui.Log().Print(fds.GetGenre(), fds)

	return
}

func (d *objects) ReadObjectsAndBlobs(f interfaces.FuncIter[*FDSet]) (err error) {
	for objectIdString, fds := range d.fds {
		if err = d.processFDSet(objectIdString, fds); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = f(fds); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
