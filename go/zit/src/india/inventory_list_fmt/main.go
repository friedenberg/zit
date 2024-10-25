package inventory_list_fmt

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_probe_index"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type VersionedFormat interface {
	WriteInventoryListBlob(*sku.List, func() (sha.WriteCloser, error)) (*sha.Sha, error)
	WriteInventoryListObject(*sku.Transacted, func() (sha.WriteCloser, error)) (*sha.Sha, error)
	ReadInventoryListObject(io.Reader) (int64, *sku.Transacted, error)
	StreamInventoryListBlobSkus(
		rf func(interfaces.ShaGetter) (interfaces.ShaReadCloser, error),
		blobSha interfaces.Sha,
		f interfaces.FuncIter[*sku.Transacted],
	) error
}

type FormatInventoryListPrinter interface {
	Offset() int64
	Print(object_inventory_format.FormatterContext) (int64, error)
	PrintMany(...object_inventory_format.FormatterContext) (int64, error)
}

type FormatInventoryListScanner interface {
	Error() error
	GetTransacted() *sku.Transacted
	GetRange() object_probe_index.Range
	Scan() bool
	SetDebug()
}
