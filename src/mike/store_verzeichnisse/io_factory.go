package store_verzeichnisse

import (
	"io"

	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/lima/zettel_verzeichnisse"
)

type readCloserFactory interface {
	ReadCloserVerzeichnisse(string) (sha.ReadCloser, error)
}

type writeCloserFactory interface {
	WriteCloserVerzeichnisse(string) (sha.WriteCloser, error)
}

type ZettelVerzeichnisseWriterGetter interface {
	ZettelVerzeichnisseWriter(int) collections.WriterFunc[*zettel_verzeichnisse.Verzeichnisse]
}

type PageHeader interface {
	PageHeaderReaderFrom(int) io.ReaderFrom
	PageHeaderWriterTo(int) io.WriterTo
}

type ioFactory interface {
	readCloserFactory
	writeCloserFactory
}
