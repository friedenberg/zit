package objekte

import (
	"os"
)

type Mover struct {
	basePath string
	file     *os.File
	*writer
}

func NewWriterMover(age _Age, basePath string) (m Mover, err error) {
	m = Mover{
		basePath: basePath,
	}

	if m.file, err = _TempFile(); err != nil {
		err = _Error(err)
		return
	}

	if m.writer, err = NewWriter(age, m.file); err != nil {
		err = _Error(err)
		return
	}

	return
}

func (m Mover) Close() (err error) {
	if err = m.writer.Close(); err != nil {
		err = _Error(err)
		return
	}

	sha := m.writer.Sha()
	p := m.file.Name()

	if err = _Close(m.file); err != nil {
		err = _Error(err)
		return
	}

	var objektePath string

	if objektePath, err = _IdMakeDirIfNecessary(sha, m.basePath); err != nil {
		err = _Error(err)
		return
	}

	if err = os.Rename(p, objektePath); err != nil {
		err = _Error(err)
		return
	}

	return
}
