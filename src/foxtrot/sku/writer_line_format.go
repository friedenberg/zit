package sku

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
)

func String(o SkuLike) (str string) {
	str = fmt.Sprintf(
		"%s %s %s %s %s",
		o.GetTai(),
		o.GetGattung(),
		o.GetKennungLike(),
		o.GetObjekteSha(),
		o.GetAkteSha(),
	)

	return
}

func MakeWriterLineFormat(
	lf *format.LineWriter,
) schnittstellen.FuncIter[SkuLike] {
	return func(o SkuLike) (err error) {
		lf.WriteFormat("%s", o)

		return
	}
}
