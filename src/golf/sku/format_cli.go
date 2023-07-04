package sku

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
)

func MakeTimePrefixWriter[T Getter](
	f schnittstellen.FuncWriterFormat[T],
) schnittstellen.FuncWriterFormat[T] {
	return func(w io.Writer, e T) (n int64, err error) {
		t := e.GetSkuLike().GetTai()

		return format.Write(
			w,
			format.MakeFormatStringRightAligned(
				"%s",
				t.Format(kennung.FormatDateTime),
			),
			format.MakeWriter(f, e),
		)
	}
}
