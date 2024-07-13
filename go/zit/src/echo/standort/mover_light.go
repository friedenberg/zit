package standort

import (
	"bytes"
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/angeboren"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

var poolBuf interfaces.Pool[bytes.Buffer, *bytes.Buffer]

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

type MoverLight struct {
	swc             sha.WriteCloser
	buf             *bytes.Buffer
	age             *age.Age
	CompressionType angeboren.CompressionType

	basePath                  string
	objektePath               string
	lockFile                  bool
	errorOnAttemptedOverwrite bool
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

	m.age = o.Age
	m.CompressionType = o.CompressionType

	m.swc = sha.MakeWriter(m.buf)

	return
}

func (m *MoverLight) Write(p []byte) (n int, err error) {
	return m.swc.Write(p)
}

func (m *MoverLight) ReadFrom(r io.Reader) (n int64, err error) {
	return m.swc.ReadFrom(r)
}

func (m *MoverLight) GetShaLike() interfaces.ShaLike {
	return m.swc.GetShaLike()
}

// TODO handle intermittent failure on duplicates
func (m *MoverLight) Close() (err error) {
	defer poolBuf.Put(m.buf)

	if err = m.swc.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := m.GetShaLike()

	if m.objektePath == "" {
		// TODO-P3 move this validation to options
		if m.basePath == "" {
			err = errors.Errorf("basepath is nil")
			return
		}

		m.objektePath = id.Path(sh, m.basePath)
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

	if f, err = files.CreateExclusiveWriteOnlyAndMaybeMakeDir(
		m.objektePath,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var w Writer

	if w, err = NewWriter(WriteOptions{
		Age:             m.age,
		CompressionType: m.CompressionType,
		Writer:          f,
	}); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = io.Copy(w, m.buf); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = w.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = f.Close(); err != nil {
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
