package age_io

import (
	"os"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/age"
)

func NewFileWriter(a age.Age, p string) (w Writer, err error) {
	var fw *os.File

	if fw, err = open_file_guard.OpenFile(p, os.O_WRONLY|os.O_EXCL|os.O_CREATE, 0644); err != nil {
		err = errors.Error(err)
		return
	}

	var ow Writer

	if ow, err = NewWriter(a, w); err != nil {
		err = errors.Error(err)
		return
	}

	w = &fileWriter{
		file:   fw,
		Writer: ow,
	}

	return
}

type fileWriter struct {
	file *os.File
	Writer
}

func (aw fileWriter) Close() (err error) {
	if aw.file == nil {
		err = errors.Errorf("nil file")
		return
	}

	if aw.Writer == nil {
		err = errors.Errorf("nil objekte reader")
		return
	}

	if err = aw.Writer.Close(); err != nil {
		err = errors.Error(err)
		return
	}

	if err = open_file_guard.Close(aw.file); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
