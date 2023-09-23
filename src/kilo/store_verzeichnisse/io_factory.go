package store_verzeichnisse

import "github.com/friedenberg/zit/src/hotel/sku"

type PageDelegate interface {
	ShouldAddVerzeichnisse(*sku.Transacted2) error
	ShouldFlushVerzeichnisse(*sku.Transacted2) error
}

type PageDelegateGetter interface {
	GetVerzeichnissePageDelegate(int) PageDelegate
}
