package inventory_list_blobs

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

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
