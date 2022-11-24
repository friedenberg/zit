package store_objekten

import (
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/metadatei_io"
	"github.com/friedenberg/zit/src/delta/standort"
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
	readCloserFactory
	writeCloserFactory
}

type substoreAccess interface {
  LockSmith
	metadatei_io.AkteIOFactory
	ioFactory
	Age() age.Age
	Standort() standort.Standort
}
