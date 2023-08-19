package store_verzeichnisse

import "github.com/friedenberg/zit/src/india/transacted"

type PageDelegate interface {
	ShouldAddVerzeichnisse(*transacted.Zettel) error
	ShouldFlushVerzeichnisse(*transacted.Zettel) error
}

type PageDelegateGetter interface {
	GetVerzeichnissePageDelegate(int) PageDelegate
}
