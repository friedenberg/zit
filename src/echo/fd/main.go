package fd

import (
	"crypto/sha256"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/delta/sha"
	"github.com/friedenberg/zit/src/delta/ts"
)

type FD struct {
	Path    string
	ModTime ts.Time
	Sha     sha.Sha
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

	fd = FileInfo(fi)

	return
}

func FileInfo(fi os.FileInfo) FD {
	return FD{
		Path:    fi.Name(),
		ModTime: ts.Tyme(fi.ModTime()),
	}
}

// TODO-P4 move this somewhere that prevents redundant reading
func (fd *FD) Set(v string) (err error) {
	fd.Path = v

	hash := sha256.New()

	var f *os.File

	if f, err = files.Open(fd.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	if _, err = io.Copy(hash, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	fd.Sha = sha.FromHash(hash)

	return
}

//TODO-P4 add formatter
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
