package store_verzeichnisse

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type readCloserFactory interface {
	ReadCloserVerzeichnisse(string) (sha.ReadCloser, error)
}

type writeCloserFactory interface {
	WriteCloserVerzeichnisse(string) (sha.WriteCloser, error)
}

type ZettelTransactedWriterGetter interface {
	ZettelTransactedWriter(int) schnittstellen.FuncIter[*zettel.Transacted]
}

type ioFactory interface {
	readCloserFactory
	writeCloserFactory
}
