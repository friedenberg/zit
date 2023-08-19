package store_verzeichnisse

import "github.com/friedenberg/zit/src/golf/sku"

type PageDelegate interface {
	ShouldAddVerzeichnisse(*sku.TransactedZettel) error
	ShouldFlushVerzeichnisse(*sku.TransactedZettel) error
}

type PageDelegateGetter interface {
	GetVerzeichnissePageDelegate(int) PageDelegate
}
