package sku_formats

import (
	"fmt"

	"github.com/friedenberg/zit/src/hotel/sku"
)

type KeyerSkuLikeUnique struct{}

func (k KeyerSkuLikeUnique) GetKey(o sku.SkuLike) string {
	if o == nil {
		return ""
	}

	return fmt.Sprintf(
		"%s %s %s %s %s",
		o.GetTai(),
		o.GetGattung(),
		o.GetKennungLike(),
		o.GetObjekteSha(),
		o.GetAkteSha(),
	)
}

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
