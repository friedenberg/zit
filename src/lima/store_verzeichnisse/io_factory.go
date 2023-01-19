package store_verzeichnisse

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

type readCloserFactory interface {
	ReadCloserVerzeichnisse(string) (sha.ReadCloser, error)
}

type writeCloserFactory interface {
	WriteCloserVerzeichnisse(string) (sha.WriteCloser, error)
}

type ZettelTransactedWriterGetter interface {
	ZettelTransactedWriter(int) collections.WriterFunc[*zettel.Transacted]
}

type PageHeader interface {
	PageHeaderReaderFrom(int) io.ReaderFrom
	PageHeaderWriterTo(int) io.WriterTo
}

type ioFactory interface {
	readCloserFactory
	writeCloserFactory
}
