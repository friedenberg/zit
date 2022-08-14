package objekte

import (
	"os"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/age"
)

type Mover struct {
	fileWriter
	basePath    string
	objektePath string
}

func NewWriterMoverPrenamed(age age.Age, objektePath string) (m *Mover, err error) {
	m = &Mover{
		objektePath: objektePath,
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
	if err = m.fileWriter.Close(); err != nil {
		err = errors.Error(err)
		return
	}

	if m.objektePath == "" {
		if m.basePath == "" {
			err = errors.Errorf("basepath is nil")
			return
		}

		sha := m.Writer.Sha()

		if m.objektePath, err = id.MakeDirIfNecessary(sha, m.basePath); err != nil {
			err = errors.Error(err)
			return
		}
	}

	p := m.file.Name()

	if err = os.Rename(p, m.objektePath); err != nil {
		err = errors.Error(err)
		return
	}

	// if err = files.SetDisallowUserChanges(m.objektePath); err != nil {
	// 	err = errors.Error(err)
	// 	return
	// }

	return
}
