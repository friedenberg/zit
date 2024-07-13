package bestandsaufnahme

import (
	"bufio"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

type format struct {
	objekteFormat object_inventory_format.Format
	options       object_inventory_format.Options
}

func MakeFormat(
	sv interfaces.StoreVersion,
	op object_inventory_format.Options,
) format {
	return format{
		objekteFormat: object_inventory_format.FormatForVersion(sv),
		options:       op,
	}
}

func (f format) ParseBlob(
	r io.Reader,
	o *InventoryList,
) (n int64, err error) {
	dec := sku_fmt.MakeFormatBestandsaufnahmeScanner(
		r,
		f.objekteFormat,
		f.options,
	)

	// dec.SetDebug()

	for dec.Scan() {
		sk := dec.GetTransacted()

		if err = o.Skus.Add(sk); err != nil {
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

func (f format) Format(w io.Writer, o *InventoryList) (n int64, err error) {
	return f.FormatParsedInventoryList(w, o)
}

func (f format) FormatParsedInventoryList(
	w io.Writer,
	o *InventoryList,
) (n int64, err error) {
	bw := bufio.NewWriter(w)
	defer errors.DeferredFlusher(&err, bw)

	fo := sku_fmt.MakeFormatBestandsaufnahmePrinter(
		bw,
		f.objekteFormat,
		f.options,
	)

	defer o.Skus.Restore()

	var n1 int64

	for {
		sk, ok := o.Skus.PopAndSave()

		if !ok {
			break
		}

		n1, err = fo.Print(sk)

		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
