package to_merge

import (
	"fmt"

	"github.com/friedenberg/zit/src/hotel/sku"
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
		o.GetKennungLike(),
		o.GetObjekteSha(),
		o.GetAkteSha(),
	)
}

func String(o *sku.Transacted) (str string) {
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
