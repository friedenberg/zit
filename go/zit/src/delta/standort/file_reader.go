package standort

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
)

func NewFileReader(o FileReadOptions) (r Reader, err error) {
	ar := objekteReader{}

	if ar.file, err = files.Open(o.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	fro := ReadOptions{
		Age:             o.Age,
		Reader:          ar.file,
		CompressionType: o.CompressionType,
	}

	if ar.Reader, err = NewReader(fro); err != nil {
		err = errors.Wrap(err)
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
		err = errors.Wrap(err)
		return
	}

	if err = files.Close(ar.file); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
