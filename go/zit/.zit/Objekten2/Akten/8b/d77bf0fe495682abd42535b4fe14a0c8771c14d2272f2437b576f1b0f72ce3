package store_browser

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

var (
	poolExternal   interfaces.Pool[sku.Transacted, *sku.Transacted]
	poolCheckedOut interfaces.Pool[sku.CheckedOut, *sku.CheckedOut]
)

func init() {
	poolExternal = pool.MakePool[sku.Transacted](
		nil,
		nil,
	)

	poolCheckedOut = pool.MakePool[sku.CheckedOut](
		nil,
		nil,
	)
}

func GetExternalPool() interfaces.Pool[sku.Transacted, *sku.Transacted] {
	return poolExternal
}

func GetCheckedOutPool() interfaces.Pool[sku.CheckedOut, *sku.CheckedOut] {
	return poolCheckedOut
}
