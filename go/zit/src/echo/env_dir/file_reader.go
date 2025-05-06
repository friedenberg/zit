package env_dir

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
)

func NewFileReader(
	options FileReadOptions,
) (readCloser interfaces.ShaReadCloser, err error) {
	objectReader := objectReader{}

	if options.Path == "-" {
		objectReader.file = os.Stdin
	} else {
		if objectReader.file, err = files.Open(options.Path); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	readOptions := ReadOptions{
		Config: MakeConfig(
			options.GetBlobCompression(),
			options.GetBlobEncryption(),
			options.GetLockInternalFiles(),
		),
		File: objectReader.file,
	}

	if objectReader.ShaReadCloser, err = NewReader(readOptions); err != nil {
		if _, err = objectReader.file.Seek(0, io.SeekStart); err != nil {
			err = errors.Wrap(err)
			return
		}

		readOptions.encryption = &age.Age{}

		if objectReader.ShaReadCloser, err = NewReader(readOptions); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	readCloser = objectReader

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
		err = errors.ErrorWithStackf("nil file")
		return
	}

	if ar.ShaReadCloser == nil {
		err = errors.ErrorWithStackf("nil object reader")
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
