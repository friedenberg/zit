package inventory_list_fmt

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
)

func (pf Factory) MakePrinter(out io.Writer) FormatInventoryListPrinter {
	return MakePrinter(
		out,
		pf.Format,
		pf.Options,
	)
}

func (pf Factory) MakeScanner(in io.Reader) FormatInventoryListScanner {
	return MakeScanner(
		in,
		pf.Format,
		pf.Options,
	)
}

type Factory struct {
	object_inventory_format.Format
	object_inventory_format.Options
}
