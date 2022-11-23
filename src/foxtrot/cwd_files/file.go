package cwd_files

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type File struct {
	Path string
	sha.Sha
}

func (ut File) String() string {
	return fmt.Sprintf("[%s %s]", ut.Path, ut.Sha)
}

func MakeFile(dir string, p string) (ut File, err error) {
	ut = File{
		Path: path.Join(dir, p),
	}

	hash := sha256.New()

	var f *os.File

	if f, err = files.Open(ut.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer files.Close(f)

	if _, err = io.Copy(hash, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	ut.Sha = sha.FromHash(hash)

	if ut.Path, err = filepath.Rel(dir, ut.Path); err != nil {
		err = errors.Wrapf(err, "%s", ut.Path)
		return
	}

	return
}
