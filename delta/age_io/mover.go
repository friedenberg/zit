package age_io

import (
	"os"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/files"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/age"
)

type Mover struct {
	fileWriter
	basePath    string
	objektePath string
	lockFile    bool
}

func NewWriterMoverPrenamed(age age.Age, objektePath string) (m *Mover, err error) {
	return NewMoverOptions(
		MoveOptions{
			Age:                      age,
			FinalPath:                objektePath,
			GenerateFinalPathFromSha: false,
		},
	)
}

func NewWriterMover(age age.Age, basePath string) (m *Mover, err error) {
	return NewMoverOptions(
		MoveOptions{
			Age:                      age,
			FinalPath:                basePath,
			GenerateFinalPathFromSha: true,
		},
	)
}

func NewMoverOptions(o MoveOptions) (m *Mover, err error) {
	m = &Mover{
		lockFile: o.LockFile,
	}

	if o.GenerateFinalPathFromSha {
		m.basePath = o.FinalPath
	} else {
		m.objektePath = o.FinalPath
	}

	if m.file, err = open_file_guard.TempFile(); err != nil {
		err = errors.Error(err)
		return
	}

	wo := WriteOptions{
		Age:    o.Age,
		UseZip: o.UseZip,
		Writer: m.file,
	}

	if m.Writer, err = NewWriterOptions(wo); err != nil {
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

	if m.lockFile {
		if err = files.SetDisallowUserChanges(m.objektePath); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}
