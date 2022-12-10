package store_verzeichnisse

import (
	"io"

	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/delta/sha"
	"github.com/friedenberg/zit/src/juliett/zettel_verzeichnisse"
)

type readCloserFactory interface {
	ReadCloserVerzeichnisse(string) (sha.ReadCloser, error)
}

type writeCloserFactory interface {
	WriteCloserVerzeichnisse(string) (sha.WriteCloser, error)
}

type ZettelVerzeichnisseWriterGetter interface {
	ZettelVerzeichnisseWriter(int) collections.WriterFunc[*zettel_verzeichnisse.Zettel]
}

type PageHeader interface {
	PageHeaderReaderFrom(int) io.ReaderFrom
	PageHeaderWriterTo(int) io.WriterTo
}

type ioFactory interface {
	readCloserFactory
	writeCloserFactory
}
