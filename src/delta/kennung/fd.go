package kennung

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/echo/ts"
)

type (
	ObjekteFDGetter interface {
		GetObjekteFD() FD
	}

	AkteFDGetter interface {
		GetAkteFD() FD
	}

	FDPairGetter interface {
		ObjekteFDGetter
		AkteFDGetter
	}

	AkteFDSetter interface {
		SetAkteFD(FD)
	}
)

type FD struct {
	// TODO-P2 make all of these private and expose as methods
	IsDir   bool
	Path    string
	ModTime ts.Time
	Sha     sha.Sha
}

func (a FD) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a FD) Equals(b FD) bool {
	if a.Path != b.Path {
		return false
	}

	if !a.ModTime.Equals(b.ModTime) {
		return false
	}

	if !a.Sha.Equals(b.Sha) {
		return false
	}

	return true
}

func FDFromPathWithAkteWriterFactory(
	p string,
	awf schnittstellen.AkteWriterFactory,
) (fd FD, err error) {
	if p == "" {
		err = errors.Errorf("empty path")
		return
	}

	if awf == nil {
		panic("schnittstellen.AkteWriterFactory is nil")
	}

	var f *os.File

	if f, err = files.OpenExclusiveReadOnly(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	var akteWriter sha.WriteCloser

	if akteWriter, err = awf.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, akteWriter)

	if _, err = io.Copy(akteWriter, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	var fi os.FileInfo

	if fi, err = f.Stat(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fd, err = FileInfo(fi, path.Dir(p)); err != nil {
		err = errors.Wrap(err)
		return
	}

	fd.Path = p
	fd.Sha = sha.Make(akteWriter.Sha())

	return
}

func FDFromPath(p string) (fd FD, err error) {
	if p == "" {
		err = errors.Errorf("nil file desriptor")
		return
	}

	var fi os.FileInfo

	if fi, err = os.Stat(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fd, err = FileInfo(fi, path.Dir(p)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func File(f *os.File) (fd FD, err error) {
	if f == nil {
		err = errors.Errorf("nil file desriptor")
		return
	}

	var fi os.FileInfo

	if fi, err = f.Stat(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fd, err = FileInfo(fi, path.Dir(f.Name())); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func FileInfo(fi os.FileInfo, dir string) (fd FD, err error) {
	fd = FD{
		IsDir:   fi.IsDir(),
		ModTime: ts.Tyme(fi.ModTime()),
	}

	if fd.Path, err = filepath.Abs(path.Join(dir, fi.Name())); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (fd *FD) Set(v string) (err error) {
	var fi os.FileInfo

	if fi, err = os.Stat(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if *fd, err = FileInfo(fi, path.Dir(v)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
	// errors.TodoP2("move this and cache")
	// hash := sha256.New()

	// var f *os.File

	// if f, err = files.Open(fd.Path); err != nil {
	// err = errors.Wrap(err)
	// return
	// }

	// defer errors.Deferred(&err, f.Close)

	// if _, err = io.Copy(hash, f); err != nil {
	// err = errors.Wrap(err)
	// return
	// }

	// fd.Sha = sha.FromHash(hash)

	// return
}

// TODO-P4 add formatter
// func (ut File) String() string {
// 	return fmt.Sprintf("[%s %s]", ut.Path, ut.Sha)
// }

func (f FD) String() string {
	p := filepath.Clean(f.Path)

	if f.IsDir {
		return p + string(filepath.Separator)
	} else {
		return p
	}
}

func (e FD) Ext() string {
	return path.Ext(e.Path)
}

func (e FD) ExtSansDot() string {
	return strings.TrimPrefix(path.Ext(e.Path), ".")
}

func (e FD) FileNameSansExt() string {
	base := filepath.Base(e.Path)
	ext := e.Ext()
	return base[:len(base)-len(ext)]
}

func (f FD) IsEmpty() bool {
	if f.Path == "" {
		return true
	}

	// if f.ModTime.IsZero() {
	// 	return true
	// }

	return false
}

func (f FD) GetIdLike() (il IdLike, err error) {
	var h Hinweis

	if h, err = f.GetHinweis(); err == nil {
		il = h
		return
	}

	errors.TodoP1("implement Typ and Etikett")

	err = errors.Errorf("not an id")

	return
}

func (f FD) GetHinweis() (h Hinweis, err error) {
	parts := strings.Split(f.Path, string(filepath.Separator))

	switch len(parts) {
	case 0:
		fallthrough

	case 1:
		err = errors.Errorf("not enough parts: %q", parts)
		return

	default:
		parts = parts[len(parts)-2:]
	case 2:
		break
	}

	p := strings.Join(parts, string(filepath.Separator))

	p1 := p
	ext := path.Ext(p)

	if len(ext) != 0 {
		p1 = p[:len(p)-len(ext)]
	}

	if err = h.Set(p1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
