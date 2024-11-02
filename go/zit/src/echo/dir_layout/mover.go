package dir_layout

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
)

type Mover struct {
	file *os.File
	Writer

	basePath                  string
	objectPath                string
	lockFile                  bool
	errorOnAttemptedOverwrite bool
}

func NewMover(o MoveOptions) (m *Mover, err error) {
	m = &Mover{
		lockFile:                  o.LockFile,
		errorOnAttemptedOverwrite: o.ErrorOnAttemptedOverwrite,
	}

	if o.GenerateFinalPathFromSha {
		m.basePath = o.FinalPath
	} else {
		m.objectPath = o.FinalPath
	}

	if m.file, err = o.FileTemp(); err != nil {
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
		err = errors.Errorf("nil object reader")
		return
	}

	if err = m.Writer.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	// var fi os.FileInfo

	// if fi, err = m.file.Stat(); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	if err = files.Close(m.file); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := m.GetShaLike()

	// log.Log().Printf(
	// 	"wrote %d bytes to %s, sha %s",
	// 	fi.Size(),
	// 	m.file.Name(),
	// 	sh,
	// )

	if m.objectPath == "" {
		// TODO-P3 move this validation to options
		if m.basePath == "" {
			err = errors.Errorf("basepath is nil")
			return
		}

		if m.objectPath, err = id.MakeDirIfNecessary(sh, m.basePath); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	p := m.file.Name()

	if err = os.Rename(p, m.objectPath); err != nil {
		if files.Exists(m.objectPath) {
			if m.errorOnAttemptedOverwrite {
				err = MakeErrAlreadyExists(sh, m.objectPath)
			} else {
				err = nil
			}

			return
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	// log.Log().Printf("moved %s to %s", p, m.objectPath)

	if m.lockFile {
		if err = files.SetDisallowUserChanges(m.objectPath); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
