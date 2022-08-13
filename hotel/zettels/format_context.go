package zettels

import (
	"io"
	"os"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/delta/objekte"
)

type akteReader struct {
	file *os.File
	objekte.Reader
}

func (ar akteReader) Close() (err error) {
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

func (zs zettels) AkteReader(sha sha.Sha) (r io.ReadCloser, err error) {
	return
}
