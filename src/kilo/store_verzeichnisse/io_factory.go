package store_verzeichnisse

import (
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type PageDelegate interface {
	ShouldAddVerzeichnisse(*zettel.Transacted) error
	ShouldFlushVerzeichnisse(*zettel.Transacted) error
}

type PageDelegateGetter interface {
	GetVerzeichnissePageDelegate(int) PageDelegate
}
