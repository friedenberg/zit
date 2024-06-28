package chrome

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
)

var (
	poolExternal   schnittstellen.Pool[External, *External]
	poolCheckedOut schnittstellen.Pool[CheckedOut, *CheckedOut]
)

func init() {
	poolExternal = pool.MakePool[External](
		nil,
		nil,
	)

	poolCheckedOut = pool.MakePool[CheckedOut](
		nil,
		nil,
	)
}

func GetExternalPool() schnittstellen.Pool[External, *External] {
	return poolExternal
}

func GetCheckedOutPool() schnittstellen.Pool[CheckedOut, *CheckedOut] {
	return poolCheckedOut
}
