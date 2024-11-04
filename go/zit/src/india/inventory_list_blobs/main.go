package inventory_list_blobs

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_probe_index"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

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

func ReadInventoryListBlob(
	vf sku.ListFormat,
	r io.Reader,
	a *sku.List,
) (err error) {
	if err = vf.StreamInventoryListBlobSkus(
		r,
		func(sk *sku.Transacted) (err error) {
			if err = a.Add(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
