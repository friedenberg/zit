package objekte

import (
	"os"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/age"
)

type Mover struct {
	basePath string
	file *os.File
	Writer
  objektePath string
}

func NewWriterMover(age age.Age, basePath string) (m *Mover, err error) {
	m = &Mover{
		basePath: basePath,
	}

	if m.file, err = open_file_guard.TempFile(); err != nil {
		err = errors.Error(err)
		return
	}

	if m.Writer, err = NewWriter(age, m.file); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (m *Mover) Close() (err error) {
	if err = m.Writer.Close(); err != nil {
		err = errors.Error(err)
		return
	}

	sha := m.Writer.Sha()
	p := m.file.Name()

	if err = open_file_guard.Close(m.file); err != nil {
		err = errors.Error(err)
		return
	}

	if m.objektePath, err = id.MakeDirIfNecessary(sha, m.basePath); err != nil {
		err = errors.Error(err)
		return
	}

	if err = os.Rename(p, m.objektePath); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
