package store_fs

import (
	"io/fs"
	"path/filepath"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
)

type fdExtSet map[string]*fd.FD

type dir struct {
	root string
	file_extensions.FileExtensions
	files map[string]fdExtSet
}

func makeDir(p string, fe file_extensions.FileExtensions) (d dir, err error) {
	d.root = p
	d.FileExtensions = fe
	d.files = make(map[string]fdExtSet)

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

			var f *fd.FD

			if f, err = fd.MakeFromPathAndDirEntry(p, de); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = d.addFD(f); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	ui.Debug().Print(d.files)

	if err = d.processGroups(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (d *dir) addFD(
	f *fd.FD,
) (err error) {
	if f.IsDir() {
		return
	}

	var rel string

	if rel, err = filepath.Rel(d.root, f.GetPath()); err != nil {
		err = errors.Wrap(err)
		return
	}

	ext := filepath.Ext(rel)

	if ext == ".conflict" {
		rel = strings.TrimSuffix(rel, ext)
	}

	ext = filepath.Ext(rel)
	key := strings.TrimSuffix(rel, ext)

	existing, ok := d.files[key]

	if !ok {
		existing = make(fdExtSet)
		d.files[key] = existing
	}

	existing[rel] = f
	d.files[key] = existing

	return
}

func (d *dir) processGroups() (err error) {
	blobs := make(map[string]*fd.FD)
	objects := make(map[string]*fd.FD)

	for objectIdString, fds := range d.files {
		clear(blobs)
		clear(objects)

		var g genres.Genre
		var oneAndOnlyObject, oneAndOnlyBlob *fd.FD

		for _, f := range fds {
			ext := f.ExtSansDot()

			var g1 genres.Genre

			switch ext {
			case d.Zettel:
				g1 = genres.Zettel

			case d.Typ:
				g1 = genres.Type

			case d.Etikett:
				g1 = genres.Tag

			case d.Kasten:
				g1 = genres.Repo

			default: // blobs
				oneAndOnlyBlob = f
				blobs[ext] = f
				continue
			}

			objects[ext] = f

			if g == genres.Unknown {
				g = g1
			} else if g == g1 {
				// duplicate
			} else {
			}

			oneAndOnlyObject = f
		}

		if g == genres.Unknown {
			// just blobs
		} else if len(blobs) > 1 {
			// invalid blobs
			err = errors.Errorf(
				"several blobs matching object id %q: %q",
				objectIdString,
				blobs,
			)
		} else if len(objects) > 1 {
			err = errors.Errorf(
				"found more than one object: %q:%s, %q:%s",
				objects,
			)

			return
		}

		ui.Debug().Print(oneAndOnlyObject, oneAndOnlyBlob)
	}

	return
}

// collect all files and directories
// match pairs
// return pairs
