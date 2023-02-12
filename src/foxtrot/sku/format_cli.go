package sku

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/echo/ts"
)

type DataIdentityGetter interface {
	GetDataIdentity() DataIdentity
}

func MakeTimePrefixWriter[T DataIdentityGetter](
	f schnittstellen.FuncWriterFormat[T],
) schnittstellen.FuncWriterFormat[T] {
	return func(w io.Writer, e T) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAligned(
				"%s ",
				e.GetDataIdentity().GetTime().Format(ts.FormatDateTime),
			),
			format.MakeWriter(f, e),
		)
	}
}
