package transaktion

import (
	"io"

	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/golf/sku"
)

type Writer struct {
	Transaktion
}

func (w Writer) WriteTo(w1 io.Writer) (n int64, err error) {
	lw := format.NewLineWriter()

	lw.WriteStringers(w.Transaktion.Time)

	w.Transaktion.Skus.Each(
		func(o sku.SkuLikePtr) (err error) {
			lw.WriteFormat("%s", o)

			return
		},
	)

	return lw.WriteTo(w1)
}
