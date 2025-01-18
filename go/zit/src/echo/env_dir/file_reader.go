package env_dir

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
)

func NewFileReader(o FileReadOptions) (r interfaces.ShaReadCloser, err error) {
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
		Config: MakeConfig(
			o.GetBlobCompression(),
			o.GetBlobEncryption(),
			o.GetLockInternalFiles(),
		),
		Reader: ar.file,
	}

	if ar.ShaReadCloser, err = NewReader(fro); err != nil {
		err = errors.Wrap(err)
		return
	}

	r = ar

	return
}

type objectReader struct {
	file *os.File
	interfaces.ShaReadCloser
}

func (r objectReader) String() string {
	return r.file.Name()
}

func (ar objectReader) Close() (err error) {
	if ar.file == nil {
		err = errors.Errorf("nil file")
		return
	}

	if ar.ShaReadCloser == nil {
		err = errors.Errorf("nil object reader")
		return
	}

	if err = ar.ShaReadCloser.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = files.Close(ar.file); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
