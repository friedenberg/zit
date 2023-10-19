package fd

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/thyme"
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
	isDir   bool
	path    string
	modTime thyme.Time
	sha     sha.Sha
	state   State
}

func (a FD) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a FD) Equals(b FD) bool {
	if a.path != b.path {
		return false
	}

	if !a.modTime.Equals(b.modTime) {
		return false
	}

	if !a.sha.Equals(b.sha) {
		return false
	}

	return true
}

func (fd *FD) Set(v string) (err error) {
	v = strings.TrimSpace(v)

	if v == "." {
		err = errors.Errorf("'.' not supported")
		return
	}

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
	p := filepath.Clean(f.path)

	if f.isDir {
		return p + string(filepath.Separator)
	} else {
		return p
	}
}

func (e FD) Ext() string {
	return path.Ext(e.path)
}

func (e FD) ExtSansDot() string {
	return strings.TrimPrefix(path.Ext(e.path), ".")
}

func (e FD) FileNameSansExt() string {
	base := filepath.Base(e.path)
	ext := e.Ext()
	return base[:len(base)-len(ext)]
}

func (f FD) IsEmpty() bool {
	if f.path == "" {
		return true
	}

	// if f.ModTime.IsZero() {
	// 	return true
	// }

	return false
}

func (fd FD) Parts() [3]string {
	return [3]string{"", "", fd.String()}
}

func (fd FD) GetPath() string {
	return fd.path
}

func (fd *FD) SetPath(p string) (err error) {
	fd.path = p
	return
}

func (fd *FD) SetPathRel(p, dir string) (err error) {
	if p, err = filepath.Rel(dir, p); err != nil {
		err = errors.Wrapf(err, "path: %q", p)
		return
	}

	return fd.SetPath(p)
}

func (fd FD) IsDir() bool {
	return fd.isDir
}

func (fd *FD) SetShaLike(v schnittstellen.ShaLike) {
	fd.sha = sha.Make(v)
}

func (fd FD) GetShaLike() schnittstellen.ShaLike {
	return fd.sha
}

func (fd FD) GetState() State {
	return fd.state
}

func (fd *FD) Reset() {
	fd.state = StateUnknown
	fd.isDir = false
	fd.path = ""
	fd.modTime.Reset()
	fd.sha.Reset()
}
