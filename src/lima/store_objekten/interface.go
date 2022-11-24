package store_objekten

import (
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/metadatei_io"
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
	metadatei_io.AkteIOFactory
	readCloserFactory
	writeCloserFactory
}
