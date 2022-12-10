package store_objekten

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
)

type LockSmith interface {
	IsAcquired() bool
}

type readCloserFactory interface {
	ReadCloserObjekten(string) (sha.ReadCloser, error)
	ReadCloserVerzeichnisse(string) (sha.ReadCloser, error)
}

type writeCloserFactory interface {
	WriteCloserObjekten(string) (sha.WriteCloser, error)
	WriteCloserVerzeichnisse(string) (sha.WriteCloser, error)
}

type ioFactory interface {
	gattung.AkteIOFactory
	readCloserFactory
	writeCloserFactory
}
