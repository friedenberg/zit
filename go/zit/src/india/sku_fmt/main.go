package sku_fmt

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type (
	CheckedOut = sku.CheckedOutFS
	Transacted = sku.Transacted
)

type KeyerSkuLikeUnique struct{}

func (k KeyerSkuLikeUnique) GetKey(o *sku.Transacted) string {
	if o == nil {
		return ""
	}

	return fmt.Sprintf(
		"%s %s %s %s %s",
		o.GetTai(),
		o.GetGattung(),
		o.GetKennung(),
		o.GetObjekteSha(),
		o.GetAkteSha(),
	)
}

func String(o *sku.Transacted) (str string) {
	str = fmt.Sprintf(
		"%s %s %s %s %s",
		o.GetTai(),
		o.GetGattung(),
		o.GetKennung(),
		o.GetObjekteSha(),
		o.GetAkteSha(),
	)

	return
}
