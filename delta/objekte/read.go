package objekte

import (
	"io"
	"os"
)

func Read(out io.Writer, age _Age, p string) (err error) {
	var f *os.File

	if f, err = _Open(p); err != nil {
		err = _Error(err)
		return
	}

	defer _Close(f)

	var r io.ReadCloser

	if r, err = NewReader(age, f); err != nil {
		err = _Error(err)
		return
	}

	defer r.Close()

	if _, err = io.Copy(out, r); err != nil {
		err = _Error(err)
		return
	}

	return
}
