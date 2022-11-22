package metadatei_io

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
)

const (
	Boundary = "---"
)

type MetadateiWriterTo interface {
	io.WriterTo
	HasMetadateiContent() bool
}

type AkteIOFactory interface {
	AkteWriter() (sha.WriteCloser, error)
	AkteReader(sha.Sha) (sha.ReadCloser, error)
}

type AkteIOFactoryFactory interface {
	AkteFactory(gattung.Gattung) AkteIOFactory
}
