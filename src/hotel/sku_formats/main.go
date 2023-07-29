package sku_formats

import (
	"fmt"

	"github.com/friedenberg/zit/src/golf/sku"
)

func String(o sku.SkuLike) (str string) {
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
