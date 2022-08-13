package objekte

import (
	"io"
	"os"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/age"
)

func Read(out io.Writer, age age.Age, p string) (err error) {
	var f *os.File

	if f, err = open_file_guard.Open(p); err != nil {
		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(f)

	var r io.ReadCloser

	if r, err = NewReader(age, f); err != nil {
		err = errors.Error(err)
		return
	}

	defer r.Close()

	if _, err = io.Copy(out, r); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
