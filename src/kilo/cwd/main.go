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
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/india/konfig"
)

type CwdFiles struct {
	akteWriterFactory schnittstellen.AkteWriterFactory
	erworben          konfig.Compiled
	dir               string
	// TODO-P4 make private
	Zettelen  schnittstellen.MutableSetPtrLike[Zettel, *Zettel]
	Typen     schnittstellen.MutableSetPtrLike[Typ, *Typ]
	Kisten    schnittstellen.MutableSetPtrLike[Kasten, *Kasten]
	Etiketten schnittstellen.MutableSetPtrLike[Etikett, *Etikett]
	// TODO-P4 make set
	UnsureAkten      schnittstellen.MutableSetPtrLike[kennung.FD, *kennung.FD]
	EmptyDirectories []kennung.FD
}

func (fs *CwdFiles) MarkUnsureAkten(fd kennung.FD) {
	fs.UnsureAkten.Add(fd)
}

func (fs CwdFiles) EachCreatableMatchable(
	m schnittstellen.FuncIter[kennung.IdLikeGetter],
) (err error) {
	todo.Parallelize()

	if err = fs.Typen.EachPtr(
		func(e *Typ) (err error) {
			return m(e)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = fs.Etiketten.EachPtr(
		func(e *Etikett) (err error) {
			return m(e)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = fs.Kisten.EachPtr(
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
	if collections.Len(
		fs.Zettelen,
		fs.Typen,
		fs.Kisten,
		fs.Etiketten,
		fs.UnsureAkten,
	) == 0 {
		return
	}

	sb := &strings.Builder{}
	sb.WriteString(kennung.QueryGroupOpenOperator)

	hasOne := false

	writeOneIfNecessary := func(v schnittstellen.Stringer) (err error) {
		if hasOne {
			sb.WriteString(kennung.QueryOrOperator)
		}

		sb.WriteString(v.String())

		hasOne = true

		return
	}

	fs.Zettelen.Each(
		func(z Zettel) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.Typen.Each(
		func(z Typ) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.Etiketten.Each(
		func(z Etikett) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.Kisten.Each(
		func(z Kasten) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.UnsureAkten.Each(
		func(z kennung.FD) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	sb.WriteString(kennung.QueryGroupCloseOperator)

	out = sb.String()
	return
}

func (fs CwdFiles) ContainsMatchable(m kennung.Matchable) bool {
	g := gattung.Must(m)

	switch g {
	case gattung.Zettel:
		return fs.Zettelen.ContainsKey(m.GetKennungLike().String())

	case gattung.Typ:
		return fs.Typen.ContainsKey(m.GetKennungLike().String())

	case gattung.Etikett:
		return fs.Etiketten.ContainsKey(m.GetKennungLike().String())

	case gattung.Kasten:
		return fs.Kisten.ContainsKey(m.GetKennungLike().String())
	}

	return true
}

func (fs CwdFiles) GetFDs() schnittstellen.SetLike[kennung.FD] {
	fds := kennung.MakeMutableFDSet()

	kennung.FDSetAddPairs[Zettel](fs.Zettelen, fds)
	kennung.FDSetAddPairs[Typ](fs.Typen, fds)
	kennung.FDSetAddPairs[Etikett](fs.Etiketten, fds)

	return fds
}

func (fs CwdFiles) GetZettel(
	h *kennung.Hinweis,
) (z *Zettel, ok bool) {
	z, ok = fs.Zettelen.GetPtr(h.String())
	return
}

func (fs CwdFiles) GetKasten(
	h *kennung.Kasten,
) (z *Kasten, ok bool) {
	z, ok = fs.Kisten.GetPtr(h.String())
	return
}

func (fs CwdFiles) GetEtikett(
	k *kennung.Etikett,
) (e *Etikett, ok bool) {
	e, ok = fs.Etiketten.GetPtr(k.String())
	return
}

func (fs CwdFiles) GetTyp(
	k *kennung.Typ,
) (t *Typ, ok bool) {
	t, ok = fs.Typen.GetPtr(k.String())
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

	iter.ErrorWaitGroupApply[Kasten](
		wg,
		fs.Kisten,
		func(e Kasten) (err error) {
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

func makeCwdFiles(
	erworben konfig.Compiled,
	dir string,
	awf schnittstellen.AkteWriterFactory,
) (fs CwdFiles) {
	fs = CwdFiles{
		akteWriterFactory: awf,
		erworben:          erworben,
		dir:               dir,
		Kisten: collections_ptr.MakeMutableValueSet[Kasten, *Kasten](
			nil,
		),
		Typen: collections_ptr.MakeMutableValueSet[Typ, *Typ](nil),
		Zettelen: collections_ptr.MakeMutableValueSet[Zettel, *Zettel](
			nil,
		),
		Etiketten: collections_ptr.MakeMutableValueSet[Etikett, *Etikett](
			nil,
		),
		UnsureAkten: collections_ptr.MakeMutableValueSet[kennung.FD, *kennung.FD](
			nil,
		),
		EmptyDirectories: make([]kennung.FD, 0),
	}

	return
}

func MakeCwdFilesAll(
	k konfig.Compiled,
	dir string,
	awf schnittstellen.AkteWriterFactory,
) (fs CwdFiles, err error) {
	fs = makeCwdFiles(k, dir, awf)
	err = fs.readAll()
	return
}

func MakeCwdFilesExactly(
	k konfig.Compiled,
	dir string,
	awf schnittstellen.AkteWriterFactory,
	files ...string,
) (fs CwdFiles, err error) {
	fs = makeCwdFiles(k, dir, awf)
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

		var fd kennung.FD

		if fd, err = kennung.FileInfo(fi, fs.dir); err != nil {
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

func (c CwdFiles) MatcherLen() int {
	return collections.Len(
		c.Zettelen,
		c.Typen,
		c.Kisten,
		c.Etiketten,
	)
}

func (_ CwdFiles) Each(_ schnittstellen.FuncIter[kennung.Matcher]) error {
	return nil
}

func (c CwdFiles) Len() int {
	return collections.Len(
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
		var ut kennung.FD

		if ut, err = MakeFile(fs.dir, a, fs.akteWriterFactory); err != nil {
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
