package transaktion

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/line_format"
	"github.com/friedenberg/zit/src/echo/sku"
)

type Writer struct {
	Transaktion
}

func (w Writer) WriteTo(w1 io.Writer) (n int64, err error) {
	lw := line_format.NewWriter()

	lw.WriteStringers(w.Transaktion.Time)

	w.Transaktion.Each(sku.MakeWriterLineFormat(lw))

	return lw.WriteTo(w1)
}
