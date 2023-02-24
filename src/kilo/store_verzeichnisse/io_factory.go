package store_verzeichnisse

import (
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type readCloserFactory interface {
	ReadCloserVerzeichnisse(string) (sha.ReadCloser, error)
}

type writeCloserFactory interface {
	WriteCloserVerzeichnisse(string) (sha.WriteCloser, error)
}

type PageDelegate interface {
	ShouldAddVerzeichnisse(*zettel.Transacted) error
	ShouldFlushVerzeichnisse(*zettel.Transacted) error
}

type PageDelegateGetter interface {
	GetVerzeichnissePageDelegate(int) PageDelegate
}

type ioFactory interface {
	readCloserFactory
	writeCloserFactory
}
