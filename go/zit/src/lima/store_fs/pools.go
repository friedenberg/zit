package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

var (
	poolExternal   interfaces.Pool[sku.External, *sku.External]
	poolCheckedOut interfaces.Pool[CheckedOut, *CheckedOut]
)

func init() {
	poolExternal = pool.MakePool[sku.External](
		nil,
		nil,
	)

	poolCheckedOut = pool.MakePool[CheckedOut](
		nil,
		nil,
	)
}

func GetExternalPool() interfaces.Pool[sku.External, *sku.External] {
	return poolExternal
}

func GetCheckedOutPool() interfaces.Pool[CheckedOut, *CheckedOut] {
	return poolCheckedOut
}
