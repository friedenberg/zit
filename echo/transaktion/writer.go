package transaktion

import (
	"io"

	"github.com/friedenberg/zit/charlie/line_format"
)

type Writer struct {
	Transaktion
}

func (w Writer) WriteTo(w1 io.Writer) (n int64, err error) {
	lw := line_format.NewWriter()

	lw.WriteStringers(w.Transaktion.Time)

	for _, o := range w.Transaktion.Objekten {
		lw.WriteFormat(
			"%s %s %s %s %s",
			o.Type,
			o.Mutter[0],
			o.Mutter[1],
			o.Id,
			o.Sha,
		)
	}

	return lw.WriteTo(w1)
}
