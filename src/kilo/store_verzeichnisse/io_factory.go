package store_verzeichnisse

import "github.com/friedenberg/zit/src/hotel/sku"

type PageDelegate interface {
	ShouldAddVerzeichnisse(*sku.Transacted) error
	ShouldFlushVerzeichnisse(*sku.Transacted) error
}

type PageDelegateGetter interface {
	GetVerzeichnissePageDelegate(int) PageDelegate
}
