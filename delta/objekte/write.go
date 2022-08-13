package objekte

import (
	"fmt"
	"io"
	"os"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/age"
)

func Write(in io.Reader, age age.Age, basePath string) (sha sha.Sha, path string, err error) {
	var wFile *os.File

	if wFile, err = open_file_guard.TempFile(); err != nil {
		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(wFile)

	path = wFile.Name()

	var w *writer

	if w, err = NewWriter(age, wFile); err != nil {
		err = errors.Error(err)
		return
	}

	defer w.Close()

	if _, err = io.Copy(w, in); err != nil {
		err = errors.Error(err)
		return
	}

	sha = w.Sha()

	return
}

func Move(basePath string, p string, sha sha.Sha, kind fmt.Stringer) (objektePath string, err error) {
	if objektePath, err = id.MakeDirIfNecessary(sha, basePath, "Objekte", kind.String()); err != nil {
		err = errors.Error(err)
		return
	}

	if err = os.Rename(p, objektePath); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func WriteAndMove(in io.Reader, age age.Age, basePath string, kind fmt.Stringer) (sha sha.Sha, err error) {
	var p string

	if sha, p, err = Write(in, age, basePath); err != nil {
		err = errors.Error(err)
		return
	}

	if _, err = Move(basePath, p, sha, kind); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
