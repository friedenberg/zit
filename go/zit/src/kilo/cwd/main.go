package cwd

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/echo/zittish"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/konfig"
)

// TODO support globs and ignores
type CwdFiles struct {
	akteWriterFactory schnittstellen.AkteWriterFactory
	erworben          *konfig.Compiled
	dir               string
	zettelen          schnittstellen.MutableSetLike[*Zettel]
	unsureZettelen    schnittstellen.MutableSetLike[*Zettel]
	typen             schnittstellen.MutableSetLike[*Typ]
	kisten            schnittstellen.MutableSetLike[*Kasten]
	etiketten         schnittstellen.MutableSetLike[*Etikett]
	unsureAkten       fd.MutableSet
	emptyDirectories  fd.MutableSet
}

func (fs *CwdFiles) MarkUnsureAkten(f *fd.FD) (err error) {
	if f, err = fd.MakeFileFromFD(f, fs.akteWriterFactory); err != nil {
		err = errors.Wrapf(err, "%q", f)
		return
	}

	if err = fs.unsureAkten.Add(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (fs CwdFiles) EachCreatableMatchable(
	m schnittstellen.FuncIter[*sku.ExternalMaybe],
) (err error) {
	todo.Parallelize()

	if err = fs.typen.Each(
		func(e *Typ) (err error) {
			return m(e)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = fs.etiketten.Each(
		func(e *Etikett) (err error) {
			return m(e)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = fs.kisten.Each(
		func(e *Kasten) (err error) {
			return m(e)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (fs CwdFiles) String() (out string) {
	if iter.Len(
		fs.zettelen,
		fs.typen,
		fs.kisten,
		fs.etiketten,
		fs.unsureAkten,
	) == 0 {
		return
	}

	sb := &strings.Builder{}
	sb.WriteRune(zittish.OpGroupOpen)

	hasOne := false

	writeOneIfNecessary := func(v schnittstellen.Stringer) (err error) {
		if hasOne {
			sb.WriteRune(zittish.OpOr)
		}

		sb.WriteString(v.String())

		hasOne = true

		return
	}

	fs.zettelen.Each(
		func(z *Zettel) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.typen.Each(
		func(z *Typ) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.etiketten.Each(
		func(z *Etikett) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.kisten.Each(
		func(z *Kasten) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.unsureAkten.Each(
		func(z *fd.FD) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	sb.WriteRune(zittish.OpGroupClose)

	out = sb.String()
	return
}

func (fs CwdFiles) GetKennungForFD(fd *fd.FD) (k *kennung.Kennung2, err error) {
	k = kennung.GetKennungPool().Get()

	if err = k.SetFromPath(
		fd.String(),
		fs.erworben.FileExtensions,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (fs CwdFiles) ContainsSku(m *sku.Transacted) bool {
	g := gattung.Must(m)

	switch g {
	case gattung.Zettel:
		return fs.zettelen.ContainsKey(m.GetKennung().String())

	case gattung.Typ:
		return fs.typen.ContainsKey(m.GetKennung().String())

	case gattung.Etikett:
		return fs.etiketten.ContainsKey(m.GetKennung().String())

	case gattung.Kasten:
		return fs.kisten.ContainsKey(m.GetKennung().String())
	}

	return true
}

func (fs CwdFiles) GetCwdFDs() fd.Set {
	fds := fd.MakeMutableSet()

	fd.SetAddPairs(fs.zettelen, fds)
	fd.SetAddPairs(fs.typen, fds)
	fd.SetAddPairs(fs.etiketten, fds)
	fd.SetAddPairs(fs.unsureZettelen, fds)
	fs.unsureAkten.Each(fds.Add)

	return fds
}

func (fs CwdFiles) GetUnsureAkten() fd.Set {
	fds := fd.MakeMutableSet()
	fs.unsureAkten.Each(fds.Add)
	return fds
}

func (fs CwdFiles) GetEmptyDirectories() fd.Set {
	fds := fd.MakeMutableSet()
	fs.emptyDirectories.Each(fds.Add)
	return fds
}

func (fs CwdFiles) GetZettel(
	h *kennung.Hinweis,
) (z *Zettel, ok bool) {
	z, ok = fs.zettelen.Get(h.String())
	return
}

func (fs CwdFiles) GetKasten(
	h *kennung.Kasten,
) (z *Kasten, ok bool) {
	z, ok = fs.kisten.Get(h.String())
	return
}

func (fs CwdFiles) GetEtikett(
	k *kennung.Etikett,
) (e *Etikett, ok bool) {
	e, ok = fs.etiketten.Get(k.String())
	return
}

func (fs CwdFiles) GetTyp(
	k *kennung.Typ,
) (t *Typ, ok bool) {
	t, ok = fs.typen.Get(k.String())
	return
}

func (fs CwdFiles) Get(
	k schnittstellen.StringerGattungGetter,
) (t *sku.ExternalMaybe, ok bool) {
	g := gattung.Must(k.GetGattung())

	switch g {
	case gattung.Kasten:
		return fs.kisten.Get(k.String())

	case gattung.Zettel:
		return fs.zettelen.Get(k.String())

	case gattung.Typ:
		return fs.typen.Get(k.String())

	case gattung.Etikett:
		return fs.etiketten.Get(k.String())

	case gattung.Konfig:
		// TODO-P3
		return

	default:
		panic(gattung.MakeErrUnsupportedGattung(g))
	}
}

func (fs CwdFiles) All(
	f schnittstellen.FuncIter[*sku.ExternalMaybe],
) (err error) {
	wg := iter.MakeErrorWaitGroupParallel()

	iter.ErrorWaitGroupApply(
		wg,
		fs.zettelen,
		func(e *Zettel) (err error) {
			return f(e)
		},
	)

	iter.ErrorWaitGroupApply(
		wg,
		fs.typen,
		func(e *Typ) (err error) {
			return f(e)
		},
	)

	iter.ErrorWaitGroupApply(
		wg,
		fs.kisten,
		func(e *Kasten) (err error) {
			return f(e)
		},
	)

	iter.ErrorWaitGroupApply(
		wg,
		fs.etiketten,
		func(e *Etikett) (err error) {
			return f(e)
		},
	)

	return wg.GetError()
}

func (fs CwdFiles) AllUnsure(
	f schnittstellen.FuncIter[*sku.ExternalMaybe],
) (err error) {
	wg := iter.MakeErrorWaitGroupParallel()

	iter.ErrorWaitGroupApply(
		wg,
		fs.unsureZettelen,
		func(e *Zettel) (err error) {
			return f(e)
		},
	)

	return wg.GetError()
}

func (fs CwdFiles) ZettelFiles() (out []string, err error) {
	out, err = iter.DerivedValues(
		fs.zettelen,
		func(z *Zettel) (p string, err error) {
			p = z.GetObjekteFD().GetPath()
			return
		},
	)

	return
}

func makeCwdFiles(
	erworben *konfig.Compiled,
	st standort.Standort,
) (fs *CwdFiles) {
	fs = &CwdFiles{
		akteWriterFactory: st,
		erworben:          erworben,
		dir:               st.Cwd(),
		kisten: collections_value.MakeMutableValueSet[*Kasten](
			nil,
		),
		typen: collections_value.MakeMutableValueSet[*Typ](nil),
		zettelen: collections_value.MakeMutableValueSet[*Zettel](
			nil,
		),
		unsureZettelen: collections_value.MakeMutableValueSet[*Zettel](
			nil,
		),
		etiketten: collections_value.MakeMutableValueSet[*Etikett](
			nil,
		),
		unsureAkten: collections_value.MakeMutableValueSet[*fd.FD](
			nil,
		),
		emptyDirectories: collections_value.MakeMutableValueSet[*fd.FD](
			nil,
		),
	}

	return
}

func MakeCwdFilesAll(
	k *konfig.Compiled,
	st standort.Standort,
) (fs *CwdFiles, err error) {
	fs = makeCwdFiles(k, st)
	err = fs.readAll()
	return
}

func MakeCwdFilesExactly(
	k *konfig.Compiled,
	st standort.Standort,
	files ...string,
) (fs *CwdFiles, err error) {
	fs = makeCwdFiles(k, st)
	err = fs.readInputFiles(files...)
	return
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

func (cwd *CwdFiles) readAll() (err error) {
	// TODO use walkdir instead
	// check for empty directories
	if err = filepath.WalkDir(
		cwd.dir,
		func(p string, d fs.DirEntry, in error) (err error) {
			if in != nil {
				err = errors.Wrap(in)
				return
			}

			var rel string

			if rel, err = filepath.Rel(cwd.dir, p); err != nil {
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

	if dirs, err = files.ReadDirNames(cwd.dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, d := range dirs {
		if strings.HasPrefix(d, ".") {
			continue
		}

		d2 := path.Join(cwd.dir, d)

		var fi os.FileInfo

		if fi, err = os.Stat(d); err != nil {
			err = errors.Wrap(err)
			return
		}

		var f *fd.FD

		if f, err = fd.FileInfo(fi, cwd.dir); err != nil {
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
				if err = cwd.emptyDirectories.Add(f); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			for _, a := range dirs2 {
				if err = cwd.readSecondLevelFile(d2, a); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

		} else if fi.Mode().IsRegular() {
			if err = cwd.readNotSecondLevelFile(d); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

func (c CwdFiles) MatcherLen() int {
	return iter.Len(
		c.zettelen,
		c.typen,
		c.kisten,
		c.etiketten,
	)
}

func (CwdFiles) Each(_ schnittstellen.FuncIter[sku.Query]) error {
	return nil
}

func (c CwdFiles) Len() int {
	return iter.Len(
		c.zettelen,
		c.typen,
		c.kisten,
		c.etiketten,
	)
}

func (fs *CwdFiles) readNotSecondLevelFile(name string) (err error) {
	if strings.HasPrefix(name, ".") {
		return
	}

	fullPath := path.Join(fs.dir, name)

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

	switch strings.TrimPrefix(ext, ".") {
	case fs.erworben.FileExtensions.Etikett:
		if err = fs.tryEtikett(fi, fs.dir); err != nil {
			err = errors.Wrap(err)
			return
		}

	case fs.erworben.FileExtensions.Kasten:
		if err = fs.tryKasten(fi, fs.dir); err != nil {
			err = errors.Wrap(err)
			return
		}

	case fs.erworben.FileExtensions.Typ:
		if err = fs.tryTyp(fi, fs.dir); err != nil {
			err = errors.Wrap(err)
			return
		}

	case fs.erworben.FileExtensions.Zettel:
		if err = fs.tryZettel(fs.dir, name, fullPath, true); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		if err = fs.addUnsureAkten("", fullPath); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (fs *CwdFiles) addUnsureAkten(dir, name string) (err error) {
	var ut *fd.FD

	fullPath := name

	if dir != "" {
		fullPath = path.Join(dir, fullPath)
	}

	if ut, err = fd.MakeFile(
		fullPath,
		fs.akteWriterFactory,
	); err != nil {
		err = errors.Wrapf(err, "Dir: %q, Name: %q", dir, name)
		return
	}

	err = fs.unsureAkten.Add(ut)

	return
}

func (fs *CwdFiles) readSecondLevelFile(dir string, name string) (err error) {
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

	ext := filepath.Ext(fullPath)
	ext = strings.ToLower(ext)
	ext = strings.TrimSpace(ext)

	switch strings.TrimPrefix(ext, ".") {
	case fs.erworben.FileExtensions.Zettel:
		fallthrough

		// Zettel-Akten can have any extension, and so default is Zettel
	default:
		if err = fs.tryZettel(dir, name, fullPath, false); err != nil {
			if err = fs.addUnsureAkten(dir, name); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}
