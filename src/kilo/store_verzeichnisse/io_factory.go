package store_verzeichnisse

import "github.com/friedenberg/zit/src/hotel/transacted"

type PageDelegate interface {
	ShouldAddVerzeichnisse(*transacted.Zettel) error
	ShouldFlushVerzeichnisse(*transacted.Zettel) error
}

type PageDelegateGetter interface {
	GetVerzeichnissePageDelegate(int) PageDelegate
}
