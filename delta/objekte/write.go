package objekte

import (
	"fmt"
	"io"
	"os"
)

func Write(in io.Reader, age _Age, basePath string) (sha _Sha, path string, err error) {
	var wFile *os.File

	if wFile, err = _TempFile(); err != nil {
		err = _Error(err)
		return
	}

	defer _Close(wFile)

	path = wFile.Name()

	var w *writer

	if w, err = NewWriter(age, wFile); err != nil {
		err = _Error(err)
		return
	}

	defer w.Close()

	if _, err = io.Copy(w, in); err != nil {
		err = _Error(err)
		return
	}

	sha = w.Sha()

	return
}

func Move(basePath string, p string, sha _Sha, kind fmt.Stringer) (objektePath string, err error) {
	if objektePath, err = _IdMakeDirIfNecessary(sha, basePath, "Objekte", kind.String()); err != nil {
		err = _Error(err)
		return
	}

	if err = os.Rename(p, objektePath); err != nil {
		err = _Error(err)
		return
	}

	return
}

func WriteAndMove(in io.Reader, age _Age, basePath string, kind fmt.Stringer) (sha _Sha, err error) {
	var p string

	if sha, p, err = Write(in, age, basePath); err != nil {
		err = _Error(err)
		return
	}

	if _, err = Move(basePath, p, sha, kind); err != nil {
		err = _Error(err)
		return
	}

	return
}
