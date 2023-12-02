package standort

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/bravo/log"
)

type Mover struct {
	file *os.File
	Writer

	basePath                  string
	objektePath               string
	lockFile                  bool
	errorOnAttemptedOverwrite bool
}

func (s Standort) NewMover(o MoveOptions) (m *Mover, err error) {
	m = &Mover{
		lockFile:                  o.LockFile,
		errorOnAttemptedOverwrite: o.ErrorOnAttemptedOverwrite,
	}

	if o.GenerateFinalPathFromSha {
		m.basePath = o.FinalPath
	} else {
		m.objektePath = o.FinalPath
	}

	if m.file, err = s.FileTempLocal(); err != nil {
		err = errors.Wrap(err)
		return
	}

	wo := WriteOptions{
		Age:             o.Age,
		CompressionType: o.CompressionType,
		Writer:          m.file,
	}

	if m.Writer, err = NewWriter(wo); err != nil {
		err = errors.Wrap(err)
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
		err = errors.Wrap(err)
		return
	}

	var fi os.FileInfo

	if fi, err = m.file.Stat(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = files.Close(m.file); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := m.Writer.GetShaLike()

	log.Log().Printf(
		"wrote %d bytes to %s, sha %s",
		fi.Size(),
		m.file.Name(),
		sh,
	)

	if m.objektePath == "" {
		// TODO-P3 move this validation to options
		if m.basePath == "" {
			err = errors.Errorf("basepath is nil")
			return
		}

		if m.objektePath, err = id.MakeDirIfNecessary(sh, m.basePath); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	p := m.file.Name()

	if err = os.Rename(p, m.objektePath); err != nil {
		if files.Exists(m.objektePath) {
			if m.errorOnAttemptedOverwrite {
				err = MakeErrAlreadyExists(sh, m.objektePath)
			} else {
				err = nil
			}

			return
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	log.Log().Printf("moved %s to %s", p, m.objektePath)

	if m.lockFile {
		if err = files.SetDisallowUserChanges(m.objektePath); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
