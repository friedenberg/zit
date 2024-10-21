package inventory_list_fax

import "io"

func MakePrinterFactory() PrinterFactory {
	return PrinterFactory{}
}

func (pf PrinterFactory) Make(out io.Writer) FormatInventoryListPrinter {
	return MakePrinter(
		out,
		pf.printer.format,
		pf.printer.options,
	)
}

type PrinterFactory struct {
	printer
}
