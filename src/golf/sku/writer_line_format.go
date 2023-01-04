package sku

import (
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/delta/format"
)

func MakeWriterLineFormat(
	lf *format.Writer,
) collections.WriterFunc[SkuLike] {
	return func(o SkuLike) (err error) {
		lf.WriteFormat(
			"%s %s %s %s %s",
			o.GetGattung(),
			o.GetMutter()[0],
			o.GetMutter()[1],
			o.GetId(),
			o.GetObjekteSha(),
		)

		return
	}
}
