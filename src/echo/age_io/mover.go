package age_io

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/files"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/delta/id"
)

type Mover struct {
	file *os.File
	Writer

	basePath    string
	objektePath string
	lockFile    bool
}

func NewMover(o MoveOptions) (m *Mover, err error) {
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

	if m.Writer, err = NewWriter(wo); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (m *Mover) Close() (err error) {
	if m.file == nil {
		err = errors.Errorf("nil file")
		return
	}

	if m.Writer == nil {
		err = errors.Errorf("nil objekte reader")
		return
	}

	if err = m.Writer.Close(); err != nil {
		err = errors.Error(err)
		return
	}

	if err = open_file_guard.Close(m.file); err != nil {
		err = errors.Error(err)
		return
	}

	sha := m.Writer.Sha()

	if m.objektePath == "" {
		//TODO move this validation to options
		if m.basePath == "" {
			err = errors.Errorf("basepath is nil")
			return
		}

		if m.objektePath, err = id.MakeDirIfNecessary(sha, m.basePath); err != nil {
			err = errors.Error(err)
			return
		}
	}

  //TODO create options for handling already exists as an error
	if m.lockFile && files.Exists(m.objektePath) {
		err = ErrAlreadyExists{
			Sha:  sha,
			Path: m.objektePath,
		}

		logz.Print(err)
		err = nil

		return
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
