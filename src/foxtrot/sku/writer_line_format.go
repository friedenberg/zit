package sku

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
)

// func String(o SkuLike) (str string) {
// 	str = fmt.Sprintf(
// 		"%s %s %s %s %s",
// 		o.GetGattung(),
// 		o.GetMutter()[0],
// 		o.GetMutter()[1],
// 		o.GetId(),
// 		o.GetObjekteSha(),
// 	)

// 	return
// }

func MakeWriterLineFormat(
	lf *format.LineWriter,
) schnittstellen.FuncIter[SkuLike] {
	return func(o SkuLike) (err error) {
		lf.WriteFormat("%s", o)

		return
	}
}
