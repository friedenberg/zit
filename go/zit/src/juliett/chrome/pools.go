package chrome

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
)

var (
	poolExternal   interfaces.Pool[External, *External]
	poolCheckedOut interfaces.Pool[CheckedOut, *CheckedOut]
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

func GetExternalPool() interfaces.Pool[External, *External] {
	return poolExternal
}

func GetCheckedOutPool() interfaces.Pool[CheckedOut, *CheckedOut] {
	return poolCheckedOut
}
