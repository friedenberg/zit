package sku

import (
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/line_format"
)

func MakeWriterLineFormat(
	lf *line_format.Writer,
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
