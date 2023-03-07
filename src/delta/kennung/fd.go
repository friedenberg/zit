package kennung

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/echo/ts"
)

type ObjekteFDGetter interface {
	GetObjekteFD() FD
}

type AkteFDGetter interface {
	GetAkteFD() FD
}

type FDPairGetter interface {
	ObjekteFDGetter
	AkteFDGetter
}

type FD struct {
	// TODO make all of these private and expose as methods
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

	if fd, err = FileInfo(fi); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func FileInfo(fi os.FileInfo) (fd FD, err error) {
	fd = FD{
		IsDir:   fi.IsDir(),
		Path:    fi.Name(),
		ModTime: ts.Tyme(fi.ModTime()),
	}

	if fd.Path, err = filepath.Abs(fd.Path); err != nil {
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

	if *fd, err = FileInfo(fi); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
	// errors.TodoP0("move this and cache")
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
	return f.Path
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
