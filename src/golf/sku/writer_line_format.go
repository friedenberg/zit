package sku

import (
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/delta/format"
)

func MakeWriterLineFormat(
	lf *format.Writer,
) collections.WriterFunc[*Sku] {
	return func(o *Sku) (err error) {
		lf.WriteFormat(
			"%s %s %s %s %s",
			o.Gattung,
			o.Mutter[0],
			o.Mutter[1],
			o.Id,
			o.Sha,
		)

		return
	}
}