package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/india/konfig"
)

type LockSmith interface {
	IsAcquired() bool
}

type ioFactory interface {
	konfig.Getter
	schnittstellen.AkteIOFactory
	//TODO-P4 move to Standort
	ReadCloserVerzeichnisse(string) (sha.ReadCloser, error)
	WriteCloserVerzeichnisse(string) (sha.WriteCloser, error)
}
