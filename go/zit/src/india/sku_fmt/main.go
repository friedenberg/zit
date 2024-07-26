package sku_fmt

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func String(o *sku.Transacted) (str string) {
	str = fmt.Sprintf(
		"%s %s %s %s %s",
		o.GetTai(),
		o.GetGenre(),
		o.GetObjectId(),
		o.GetObjectSha(),
		o.GetBlobSha(),
	)

	return
}
