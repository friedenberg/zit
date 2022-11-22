package metadatei_io

import (
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type AkteIOFactory interface {
	AkteReaderFactory
	AkteWriterFactory
}

type AkteReaderFactory interface {
	AkteReader(sha.Sha) (sha.ReadCloser, error)
}

type AkteWriterFactory interface {
	AkteWriter() (sha.WriteCloser, error)
}

type AkteIOFactoryFactory interface {
	AkteFactory(gattung.Gattung) AkteIOFactory
}
