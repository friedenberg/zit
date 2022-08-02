package zettels

import (
	"io"
	"os"
	"path"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/delta/objekte"
	// "github.com/friedenberg/zit/alfa/stdprinter"
)

func (zs zettels) AkteWriter() (w objekte.Writer, err error) {
	return objekte.NewWriterMover(zs.age, path.Join(zs.basePath, "Objekte", "Akte"))
}

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
	ar := akteReader{}

	p := id.Path(sha, zs.basePath, "Objekte", "Akte")

	if ar.file, err = open_file_guard.Open(p); err != nil {
		err = errors.Error(err)
		return
	}

	if ar.Reader, err = objekte.NewReader(zs.age, ar.file); err != nil {
		err = errors.Error(err)
		return
	}

	r = ar

	return
}
