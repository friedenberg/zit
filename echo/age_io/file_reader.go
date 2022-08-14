package age_io

import (
	"io"
	"os"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/charlie/open_file_guard"
)

func NewFileReader(o FileReadOptions) (r io.ReadCloser, err error) {
	ar := objekteReader{}

	if ar.file, err = open_file_guard.Open(o.Path); err != nil {
		err = errors.Error(err)
		return
	}

	if ar.Reader, err = NewReader(ReadOptions{Age: o.Age, Reader: ar.file}); err != nil {
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

	if err = open_file_guard.Close(ar.file); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
