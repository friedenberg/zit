package sku

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/echo/ts"
)

func MakeTimePrefixWriter[T DataIdentityGetter](
	f schnittstellen.FuncWriterFormat[T],
) schnittstellen.FuncWriterFormat[T] {
	return func(w io.Writer, e T) (n int64, err error) {
		t := e.GetDataIdentity().GetTai()

		return format.Write(
			w,
			format.MakeFormatStringRightAligned(
				"%s",
				t.Format(ts.FormatDateTime),
			),
			format.MakeWriter(f, e),
		)
	}
}
