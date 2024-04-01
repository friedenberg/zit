package cwd

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/src/charlie/files"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/echo/standort"
	"code.linenisgreat.com/zit/src/echo/zittish"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/juliett/konfig"
)

type CwdFiles struct {
	akteWriterFactory schnittstellen.AkteWriterFactory
	erworben          *konfig.Compiled
	dir               string
	// TODO-P4 make private
	Zettelen  schnittstellen.MutableSetLike[*Zettel]
	Typen     schnittstellen.MutableSetLike[*Typ]
	Kisten    schnittstellen.MutableSetLike[*Kasten]
	Etiketten schnittstellen.MutableSetLike[*Etikett]
	// TODO-P4 make set
	UnsureAkten      fd.MutableSet
	EmptyDirectories []*fd.FD
}

func (fs *CwdFiles) MarkUnsureAkten(f *fd.FD) (err error) {
	if f, err = fd.MakeFileFromFD(f, fs.akteWriterFactory); err != nil {
		err = errors.Wrapf(err, "%q", f)
		return
	}

	if err = fs.UnsureAkten.Add(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (fs CwdFiles) EachCreatableMatchable(
	m schnittstellen.FuncIter[*sku.ExternalMaybe],
) (err error) {
	todo.Parallelize()

	if err = fs.Typen.Each(
		func(e *Typ) (err error) {
			return m(e)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = fs.Etiketten.Each(
		func(e *Etikett) (err error) {
			return m(e)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = fs.Kisten.Each(
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
		fs.Zettelen,
		fs.Typen,
		fs.Kisten,
		fs.Etiketten,
		fs.UnsureAkten,
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

	fs.Zettelen.Each(
		func(z *Zettel) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.Typen.Each(
		func(z *Typ) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.Etiketten.Each(
		func(z *Etikett) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.Kisten.Each(
		func(z *Kasten) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.UnsureAkten.Each(
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
		return fs.Zettelen.ContainsKey(m.GetKennung().String())

	case gattung.Typ:
		return fs.Typen.ContainsKey(m.GetKennung().String())

	case gattung.Etikett:
		return fs.Etiketten.ContainsKey(m.GetKennung().String())

	case gattung.Kasten:
		return fs.Kisten.ContainsKey(m.GetKennung().String())
	}

	return true
}

func (fs CwdFiles) GetCwdFDs() fd.Set {
	fds := fd.MakeMutableSet()

	fd.SetAddPairs[*Zettel](fs.Zettelen, fds)
	fd.SetAddPairs[*Typ](fs.Typen, fds)
	fd.SetAddPairs[*Etikett](fs.Etiketten, fds)
	fs.UnsureAkten.Each(fds.Add)

	return fds
}

func (fs CwdFiles) GetZettel(
	h *kennung.Hinweis,
) (z *Zettel, ok bool) {
	z, ok = fs.Zettelen.Get(h.String())
	return
}

func (fs CwdFiles) GetKasten(
	h *kennung.Kasten,
) (z *Kasten, ok bool) {
	z, ok = fs.Kisten.Get(h.String())
	return
}

func (fs CwdFiles) GetEtikett(
	k *kennung.Etikett,
) (e *Etikett, ok bool) {
	e, ok = fs.Etiketten.Get(k.String())
	return
}

func (fs CwdFiles) GetTyp(
	k *kennung.Typ,
) (t *Typ, ok bool) {
	t, ok = fs.Typen.Get(k.String())
	return
}

func (fs CwdFiles) Get(
	k schnittstellen.StringerGattungGetter,
) (t *sku.ExternalMaybe, ok bool) {
	g := gattung.Must(k.GetGattung())

	switch g {
	case gattung.Kasten:
		return fs.Kisten.Get(k.String())

	case gattung.Zettel:
		return fs.Zettelen.Get(k.String())

	case gattung.Typ:
		return fs.Typen.Get(k.String())

	case gattung.Etikett:
		return fs.Etiketten.Get(k.String())

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

	iter.ErrorWaitGroupApply[*Zettel](
		wg,
		fs.Zettelen,
		func(e *Zettel) (err error) {
			return f(e)
		},
	)

	iter.ErrorWaitGroupApply[*Typ](
		wg,
		fs.Typen,
		func(e *Typ) (err error) {
			return f(e)
		},
	)

	iter.ErrorWaitGroupApply[*Kasten](
		wg,
		fs.Kisten,
		func(e *Kasten) (err error) {
			return f(e)
		},
	)

	iter.ErrorWaitGroupApply[*Etikett](
		wg,
		fs.Etiketten,
		func(e *Etikett) (err error) {
			return f(e)
		},
	)

	return wg.GetError()
}

func (fs CwdFiles) ZettelFiles() (out []string, err error) {
	out, err = iter.DerivedValues[*Zettel, string](
		fs.Zettelen,
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
		Kisten: collections_value.MakeMutableValueSet[*Kasten](
			nil,
		),
		Typen: collections_value.MakeMutableValueSet[*Typ](nil),
		Zettelen: collections_value.MakeMutableValueSet[*Zettel](
			nil,
		),
		Etiketten: collections_value.MakeMutableValueSet[*Etikett](
			nil,
		),
		UnsureAkten: collections_value.MakeMutableValueSet[*fd.FD](
			nil,
		),
		EmptyDirectories: make([]*fd.FD, 0),
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

		var f *fd.FD

		if f, err = fd.FileInfo(fi, fs.dir); err != nil {
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
				fs.EmptyDirectories = append(fs.EmptyDirectories, f)
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

func (c CwdFiles) MatcherLen() int {
	return iter.Len(
		c.Zettelen,
		c.Typen,
		c.Kisten,
		c.Etiketten,
	)
}

func (CwdFiles) Each(_ schnittstellen.FuncIter[sku.Query]) error {
	return nil
}

func (c CwdFiles) Len() int {
	return iter.Len(
		c.Zettelen,
		c.Typen,
		c.Kisten,
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

	default:
		var ut *fd.FD

		if ut, err = fd.MakeFile(fs.dir, a, fs.akteWriterFactory); err != nil {
			err = errors.Wrap(err)
			return
		}

		err = fs.UnsureAkten.Add(ut)
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
