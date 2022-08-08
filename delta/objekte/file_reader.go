package objekte

import (
	"io"
	"os"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/age"
)

func NewFileReader(a age.Age, p string) (r io.ReadCloser, err error) {
	ar := objekteReader{}

	if ar.file, err = open_file_guard.Open(p); err != nil {
		err = errors.Error(err)
		return
	}

	if ar.Reader, err = NewReader(a, ar.file); err != nil {
		err = errors.Error(err)
		return
	}

	r = ar

	return
}

type objekteReader struct {
	file *os.File
	Reader
}

func (ar objekteReader) Close() (err error) {
	if ar.file == nil {
		err = errors.Errorf("nil file")
		return
	}

	if ar.Reader == nil {
		err = errors.Errorf("nil objekte reader")
		return
	}

	if err = ar.Reader.Close(); err != nil {
		err = errors.Error(err)
		return
	}

	// if err = open_file_guard.Close(ar.file); err != nil {
	// 	err = errors.Error(err)
	// 	return
	// }

	return
}
