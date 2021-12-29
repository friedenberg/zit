package zettels

import (
	"io"
	"os"
	"path"
	// "github.com/friedenberg/zit/alfa/stdprinter"
)

func (zs zettels) AkteWriter() (w _ObjekteWriter, err error) {
	return _ObjekteNewWriterMover(zs.age, path.Join(zs.basePath, "Objekte", "Akte"))
}

type akteReader struct {
	file *os.File
	_ObjekteReader
}

func (ar akteReader) Close() (err error) {
	if ar.file == nil {
		err = _Errorf("nil file")
		return
	}

	if ar._ObjekteReader == nil {
		err = _Errorf("nil objekte reader")
		return
	}

	if err = ar._ObjekteReader.Close(); err != nil {
		err = _Error(err)
		return
	}

	// if err = _Close(ar.file); err != nil {
	// 	err = _Error(err)
	// 	return
	// }

	return
}

func (zs zettels) AkteReader(sha _Sha) (r io.ReadCloser, err error) {
	ar := akteReader{}

	p := _IdPath(sha, zs.basePath, "Objekte", "Akte")

	if ar.file, err = _Open(p); err != nil {
		err = _Error(err)
		return
	}

	if ar._ObjekteReader, err = _ObjekteNewReader(zs.age, ar.file); err != nil {
		err = _Error(err)
		return
	}

	r = ar

	return
}
