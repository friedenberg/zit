package inventory_list

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/india/inventory_list_fmt"
)

type format struct {
	object_inventory_format.Format
	object_inventory_format.Options
}

func MakeFormat(
	sv interfaces.StoreVersion,
	op object_inventory_format.Options,
) format {
	return format{
		Format:  object_inventory_format.FormatForVersion(sv),
		Options: op,
	}
}

func (f format) ParseBlob(
	r io.Reader,
	o *InventoryList,
) (n int64, err error) {
	dec := inventory_list_fmt.MakeScanner(
		r,
		f.Format,
		f.Options,
	)

	// dec.SetDebug()

	for dec.Scan() {
		sk := dec.GetTransacted()

		if err = o.Add(sk); err != nil {
			err = errors.Wrapf(err, "Sku: %s", sk)
			return
		}
	}

	if err = dec.Error(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
