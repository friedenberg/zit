package store_objekten

import (
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/juliett/konfig"
	"github.com/friedenberg/zit/src/schnittstellen"
)

type LockSmith interface {
	IsAcquired() bool
}

// TODO-P1 replace with standort
type readCloserFactory interface {
	ReadCloserObjekten(string) (sha.ReadCloser, error)
	ReadCloserVerzeichnisse(string) (sha.ReadCloser, error)
}

// TODO-P1 replace with standort
type writeCloserFactory interface {
	WriteCloserObjekten(string) (sha.WriteCloser, error)
	WriteCloserVerzeichnisse(string) (sha.WriteCloser, error)
}

type ioFactory interface {
	konfig.Getter
	schnittstellen.AkteIOFactory
	readCloserFactory
	writeCloserFactory
}
