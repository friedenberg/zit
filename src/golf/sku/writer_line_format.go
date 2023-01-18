package sku

import (
	"fmt"

	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/format"
)

func String(o SkuLike) (str string) {
	str = fmt.Sprintf(
		"%s %s %s %s %s",
		o.GetGattung(),
		o.GetMutter()[0],
		o.GetMutter()[1],
		o.GetId(),
		o.GetObjekteSha(),
	)

	return
}

func MakeWriterLineFormat(
	lf *format.LineWriter,
) collections.WriterFunc[SkuLike] {
	return func(o SkuLike) (err error) {
		lf.WriteFormat("%s", String(o))

		return
	}
}
