package dir_layout

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
)

func NewFileReader(o FileReadOptions) (r Reader, err error) {
	ar := objectReader{}

	if o.Path == "-" {
		ar.file = os.Stdin
	} else {
		if ar.file, err = files.Open(o.Path); err != nil {
			err = errors.Wrap(err)
			return
		}
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

type objectReader struct {
	file *os.File
	Reader
}

func (ar objectReader) Close() (err error) {
	if ar.file == nil {
		err = errors.Errorf("nil file")
		return
	}

	if ar.Reader == nil {
		err = errors.Errorf("nil object reader")
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
