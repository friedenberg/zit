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
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/echo/zittish"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func init() {
	gob.Register(External{})
}

// TODO support globs and ignores
type Store struct {
	konfig              sku.Konfig
	deletedPrinter      schnittstellen.FuncIter[*fd.FD]
	storeFuncs          sku.StoreFuncs
	metadateiTextParser metadatei.TextParser
	standort            standort.Standort
	fileEncoder         FileEncoder
	ic                  kennung.InlineTypChecker
	fileExtensions      file_extensions.FileExtensions
	dir                 string
	zettelen            schnittstellen.MutableSetLike[*KennungFDPair]
	unsureZettelen      schnittstellen.MutableSetLike[*KennungFDPair]
	typen               schnittstellen.MutableSetLike[*KennungFDPair]
	kisten              schnittstellen.MutableSetLike[*KennungFDPair]
	etiketten           schnittstellen.MutableSetLike[*KennungFDPair]
	unsureAkten         fd.MutableSet
	emptyDirectories    fd.MutableSet

	objekteFormatOptions objekte_format.Options

	deleteLock sync.Mutex
	deleted    fd.MutableSet
}

func (fs *Store) DeleteCheckout(col sku.CheckedOutLike) (err error) {
	e := col.GetSkuExternalLike().(*External)

	fs.deleteLock.Lock()
	defer fs.deleteLock.Unlock()

	if err = fs.deleted.Add(e.GetObjekteFD()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = fs.deleted.Add(e.GetAkteFD()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (fs *Store) Flush() (err error) {
	deleteOp := DeleteCheckout{}

	if err = deleteOp.Run(
		fs.konfig.IsDryRun(),
		fs.standort,
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
func (fs *Store) MarkUnsureAkten(f *fd.FD) (err error) {
	if f.IsDir() {
		// TODO handle recursion
		return
	}

	if f, err = fd.MakeFileFromFD(f, fs.standort); err != nil {
		err = errors.Wrapf(err, "%q", f)
		return
	}

	if err = fs.unsureAkten.Add(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (fs *Store) String() (out string) {
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
		func(z *KennungFDPair) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.typen.Each(
		func(z *KennungFDPair) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.etiketten.Each(
		func(z *KennungFDPair) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.kisten.Each(
		func(z *KennungFDPair) (err error) {
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

// TODO confirm against actual Kennung
func (fs *Store) GetKennungForFD(fd *fd.FD) (k *kennung.Kennung2, err error) {
	k = kennung.GetKennungPool().Get()

	if err = k.SetFromPath(
		fd.String(),
		fs.fileExtensions,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (fs *Store) ContainsSku(m *sku.Transacted) bool {
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

func (fs *Store) GetCwdFDs() fd.Set {
	fds := fd.MakeMutableSet()

	fd.SetAddPairs(fs.zettelen, fds)
	fd.SetAddPairs(fs.typen, fds)
	fd.SetAddPairs(fs.etiketten, fds)
	fd.SetAddPairs(fs.unsureZettelen, fds)
	fs.unsureAkten.Each(fds.Add)

	return fds
}

func (fs *Store) GetUnsureAkten() fd.Set {
	fds := fd.MakeMutableSet()
	fs.unsureAkten.Each(fds.Add)
	return fds
}

func (fs *Store) GetEmptyDirectories() fd.Set {
	fds := fd.MakeMutableSet()
	fs.emptyDirectories.Each(fds.Add)
	return fds
}

func (fs *Store) GetZettel(
	h *kennung.Hinweis,
) (z *KennungFDPair, ok bool) {
	z, ok = fs.zettelen.Get(h.String())
	return
}

func (fs *Store) GetKasten(
	h *kennung.Kasten,
) (z *KennungFDPair, ok bool) {
	z, ok = fs.kisten.Get(h.String())
	return
}

func (fs *Store) GetEtikett(
	k *kennung.Etikett,
) (e *KennungFDPair, ok bool) {
	e, ok = fs.etiketten.Get(k.String())
	return
}

func (fs *Store) GetTyp(
	k *kennung.Typ,
) (t *KennungFDPair, ok bool) {
	t, ok = fs.typen.Get(k.String())
	return
}

func (fs *Store) Get(
	k schnittstellen.StringerGattungGetter,
) (t *KennungFDPair, ok bool) {
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
		err := errors.Wrapf(
			gattung.MakeErrUnsupportedGattung(g),
			"Kennung: %q",
			k,
		)

		panic(err)
	}
}

func (fs *Store) All(
	f schnittstellen.FuncIter[*KennungFDPair],
) (err error) {
	wg := iter.MakeErrorWaitGroupParallel()

	iter.ErrorWaitGroupApply(
		wg,
		fs.zettelen,
		func(e *KennungFDPair) (err error) {
			return f(e)
		},
	)

	iter.ErrorWaitGroupApply(
		wg,
		fs.typen,
		func(e *KennungFDPair) (err error) {
			return f(e)
		},
	)

	iter.ErrorWaitGroupApply(
		wg,
		fs.kisten,
		func(e *KennungFDPair) (err error) {
			return f(e)
		},
	)

	iter.ErrorWaitGroupApply(
		wg,
		fs.etiketten,
		func(e *KennungFDPair) (err error) {
			return f(e)
		},
	)

	return wg.GetError()
}

func (fs *Store) AllUnsure(
	f schnittstellen.FuncIter[*KennungFDPair],
) (err error) {
	wg := iter.MakeErrorWaitGroupParallel()

	iter.ErrorWaitGroupApply(
		wg,
		fs.unsureZettelen,
		func(e *KennungFDPair) (err error) {
			return f(e)
		},
	)

	return wg.GetError()
}

func (fs *Store) ZettelFiles() (out []string, err error) {
	out, err = iter.DerivedValues(
		fs.zettelen,
		func(z *KennungFDPair) (p string, err error) {
			p = z.GetObjekteFD().GetPath()
			return
		},
	)

	return
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

func (s *Store) Initialize(ii sku.ExternalStoreInitInfo) (err error) {
  s.storeFuncs = ii.StoreFuncs
	return
}

func (s *Store) readAll() (err error) {
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

		if f, err = fd.FileInfo(fi, s.dir); err != nil {
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

func (c *Store) MatcherLen() int {
	return iter.Len(
		c.zettelen,
		c.typen,
		c.kisten,
		c.etiketten,
	)
}

func (*Store) Each(_ schnittstellen.FuncIter[sku.Query]) error {
	return nil
}

func (c *Store) Len() int {
	return iter.Len(
		c.zettelen,
		c.typen,
		c.kisten,
		c.etiketten,
	)
}

func (fs *Store) readNotSecondLevelFile(name string) (err error) {
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
	case fs.fileExtensions.Etikett:
		if err = fs.tryEtikett(fi, fs.dir); err != nil {
			err = errors.Wrap(err)
			return
		}

	case fs.fileExtensions.Kasten:
		if err = fs.tryKasten(fi, fs.dir); err != nil {
			err = errors.Wrap(err)
			return
		}

	case fs.fileExtensions.Typ:
		if err = fs.tryTyp(fi, fs.dir); err != nil {
			err = errors.Wrap(err)
			return
		}

	case fs.fileExtensions.Zettel:
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

func (fs *Store) addUnsureAkten(dir, name string) (err error) {
	var ut *fd.FD

	fullPath := name

	if dir != "" {
		fullPath = path.Join(dir, fullPath)
	}

	if ut, err = fd.MakeFile(
		fullPath,
		fs.standort,
	); err != nil {
		err = errors.Wrapf(err, "Dir: %q, Name: %q", dir, name)
		return
	}

	err = fs.unsureAkten.Add(ut)

	return
}

func (fs *Store) readSecondLevelFile(dir string, name string) (err error) {
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
	case fs.fileExtensions.Zettel:
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
