package standort

import (
	"bytes"
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/bravo/pool"
)

type MoverLight struct {
	buf *bytes.Buffer
	Writer

	basePath                  string
	objektePath               string
	lockFile                  bool
	errorOnAttemptedOverwrite bool
}

var poolBuf schnittstellen.Pool[bytes.Buffer, *bytes.Buffer]

func init() {
	poolBuf = pool.MakePool(
		func() *bytes.Buffer {
			return bytes.NewBuffer(nil)
		},
		func(b *bytes.Buffer) {
			b.Reset()
		},
	)
}

func (s Standort) NewMoverLight(o MoveOptions) (m *MoverLight, err error) {
	m = &MoverLight{
		lockFile:                  o.LockFile,
		errorOnAttemptedOverwrite: o.ErrorOnAttemptedOverwrite,
	}

	if o.GenerateFinalPathFromSha {
		m.basePath = o.FinalPath
	} else {
		m.objektePath = o.FinalPath
	}

	m.buf = poolBuf.Get()

	wo := WriteOptions{
		Age:             o.Age,
		CompressionType: o.CompressionType,
		Writer:          m.buf,
	}

	if m.Writer, err = NewWriter(wo); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (m *MoverLight) Close() (err error) {
	if m.buf == nil {
		err = errors.Errorf("nil buf")
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

	defer poolBuf.Put(m.buf)

	sh := m.GetShaLike()

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

	if files.Exists(m.objektePath) {
		if m.errorOnAttemptedOverwrite {
			err = MakeErrAlreadyExists(sh, m.objektePath)
		} else {
			err = nil
		}

		return
	}

	var f *os.File

	if f, err = files.CreateExclusiveWriteOnly(m.objektePath); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if _, err = io.Copy(f, m.buf); err != nil {
		err = errors.Wrap(err)
		return
	}

	if m.lockFile {
		if err = files.SetDisallowUserChanges(m.objektePath); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
