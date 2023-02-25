package cwd

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/india/konfig"
)

type CwdFiles struct {
	erworben         konfig.Compiled
	dir              string
	Zettelen         schnittstellen.MutableSet[Zettel]
	Typen            schnittstellen.MutableSet[Typ]
	Etiketten        schnittstellen.MutableSet[Etikett]
	UnsureAkten      []kennung.FD
	EmptyDirectories []kennung.FD
}

func (fs CwdFiles) GetMetaSet() (ms kennung.MetaSet, err error) {
	ms = kennung.MakeMetaSet(kennung.Expanders{}, gattung.Zettel)

	if err = fs.Zettelen.Each(
		func(z Zettel) (err error) {
			return ms.Add(z.Kennung, kennung.SigilNone)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = fs.Typen.Each(
		func(z Typ) (err error) {
			return ms.Add(z.Kennung, kennung.SigilNone)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = fs.Etiketten.Each(
		func(z Etikett) (err error) {
			return ms.Add(z.Kennung, kennung.SigilNone)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (fs CwdFiles) GetZettel(
	h kennung.Hinweis,
) (z Zettel, ok bool) {
	z, ok = fs.Zettelen.Get(h.String())
	return
}

func (fs CwdFiles) GetEtikett(
	k kennung.Etikett,
) (e Etikett, ok bool) {
	e, ok = fs.Etiketten.Get(k.String())
	return
}

func (fs CwdFiles) GetTyp(
	k kennung.Typ,
) (t Typ, ok bool) {
	t, ok = fs.Typen.Get(k.String())
	return
}

func (fs CwdFiles) All(
	f schnittstellen.FuncIter[sku.ExternalMaybeLike],
) (err error) {
	wg := iter.MakeErrorWaitGroup()

	iter.ErrorWaitGroupApply[Zettel](
		wg,
		fs.Zettelen,
		func(e Zettel) (err error) {
			return f(e)
		},
	)

	iter.ErrorWaitGroupApply[Typ](
		wg,
		fs.Typen,
		func(e Typ) (err error) {
			return f(e)
		},
	)

	iter.ErrorWaitGroupApply[Etikett](
		wg,
		fs.Etiketten,
		func(e Etikett) (err error) {
			return f(e)
		},
	)

	return wg.GetError()
}

func (fs CwdFiles) ZettelFiles() (out []string, err error) {
	out, err = collections.DerivedValues[Zettel, string](
		fs.Zettelen,
		func(z Zettel) (p string, err error) {
			p = z.GetObjekteFD().Path
			return
		},
	)

	return
}

func makeCwdFiles(erworben konfig.Compiled, dir string) (fs CwdFiles) {
	fs = CwdFiles{
		erworben:         erworben,
		dir:              dir,
		Typen:            collections.MakeMutableSetStringer[Typ](),
		Zettelen:         collections.MakeMutableSetStringer[Zettel](),
		Etiketten:        collections.MakeMutableSetStringer[Etikett](),
		UnsureAkten:      make([]kennung.FD, 0),
		EmptyDirectories: make([]kennung.FD, 0),
	}

	return
}

func MakeCwdFilesAll(
	k konfig.Compiled,
	dir string,
) (fs CwdFiles, err error) {
	fs = makeCwdFiles(k, dir)
	err = fs.readAll()
	return
}

func MakeCwdFilesExactly(
	k konfig.Compiled,
	dir string,
	files ...string,
) (fs CwdFiles, err error) {
	fs = makeCwdFiles(k, dir)
	err = fs.readInputFiles(files...)
	return
}

func MakeCwdFilesMetaSet(
	k konfig.Compiled,
	dir string,
	ms kennung.MetaSet,
) (fs CwdFiles, err error) {
	isZettel, ok := ms.Get(gattung.Zettel)

	switch {
	case ok && isZettel.Sigil.IncludesCwd() && isZettel.Len() == 0:
		return MakeCwdFilesAll(k, dir)

	default:
		fds := ms.GetFDs()
		files := make([]string, 0, fds.Len())

		fds.Each(
			func(fd kennung.FD) (err error) {
				files = append(files, fd.String())
				return
			},
		)

		return MakeCwdFilesExactly(k, dir, files...)
	}
}

func (fs *CwdFiles) readInputFiles(args ...string) (err error) {
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
			if err = fs.readFirstLevelFile(parts[0]); err != nil {
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

func (fs *CwdFiles) readAll() (err error) {
	var dirs []string

	if dirs, err = files.ReadDirNames(fs.dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, d := range dirs {
		if strings.HasPrefix(d, ".") {
			continue
		}

		d2 := path.Join(fs.dir, d)

		var fi os.FileInfo

		if fi, err = os.Stat(d); err != nil {
			err = errors.Wrap(err)
			return
		}

		fd := kennung.FileInfo(fi)

		if fi.Mode().IsDir() {
			var dirs2 []string

			if dirs2, err = files.ReadDirNames(d2); err != nil {
				err = errors.Wrap(err)
				return
			}

			if len(dirs2) == 0 {
				fs.EmptyDirectories = append(fs.EmptyDirectories, fd)
			}

			for _, a := range dirs2 {
				if err = fs.readSecondLevelFile(d2, a); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

		} else if fi.Mode().IsRegular() {
			if err = fs.readFirstLevelFile(d); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

func (c CwdFiles) Len() int {
	return collections.Len(
		c.Zettelen,
		c.Typen,
		c.Etiketten,
	)
}

func (fs *CwdFiles) readFirstLevelFile(a string) (err error) {
	if strings.HasPrefix(a, ".") {
		return
	}

	var fi os.FileInfo

	if fi, err = os.Stat(path.Join(fs.dir, a)); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !fi.Mode().IsRegular() {
		return
	}

	ext := filepath.Ext(a)
	ext = strings.ToLower(ext)
	ext = strings.TrimSpace(ext)

	switch strings.TrimPrefix(ext, ".") {
	case fs.erworben.FileExtensions.Etikett:
		if err = fs.tryEtikett(fi); err != nil {
			err = errors.Wrap(err)
			return
		}

	case fs.erworben.FileExtensions.Typ:
		if err = fs.tryTyp(fi); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		var ut kennung.FD

		if ut, err = MakeFile(fs.dir, a); err != nil {
			err = errors.Wrap(err)
			return
		}

		fs.UnsureAkten = append(fs.UnsureAkten, ut)
	}

	return
}

func (fs *CwdFiles) readSecondLevelFile(d string, a string) (err error) {
	if strings.HasPrefix(a, ".") {
		return
	}

	var fi os.FileInfo

	p := path.Join(d, a)

	if fi, err = os.Stat(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !fi.Mode().IsRegular() {
		return
	}

	ext := filepath.Ext(p)
	ext = strings.ToLower(ext)
	ext = strings.TrimSpace(ext)

	switch strings.TrimPrefix(ext, ".") {
	case fs.erworben.FileExtensions.Typ:
		if err = fs.tryTyp(fi); err != nil {
			err = errors.Wrap(err)
			return
		}

	case fs.erworben.FileExtensions.Zettel:
		fallthrough

		// Zettel-Akten can have any extension, and so default is Zettel
	default:
		if err = fs.tryZettel(d, a, p); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
