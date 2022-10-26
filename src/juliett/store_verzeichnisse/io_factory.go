package store_verzeichnisse

import (
	"io"

	"github.com/friedenberg/zit/src/india/zettel_verzeichnisse"
)

type readCloserFactory interface {
	ReadCloserVerzeichnisse(string) (io.ReadCloser, error)
}

type writeCloserFactory interface {
	WriteCloserVerzeichnisse(string) (io.WriteCloser, error)
}

type ZettelVerzeichnisseWriterGetter interface {
	ZettelVerzeichnisseWriter(int) zettel_verzeichnisse.Writer
}

type PageHeader interface {
	PageHeaderReaderFrom(int) io.ReaderFrom
	PageHeaderWriterTo(int) io.WriterTo
}

type ioFactory interface {
	readCloserFactory
	writeCloserFactory
}
