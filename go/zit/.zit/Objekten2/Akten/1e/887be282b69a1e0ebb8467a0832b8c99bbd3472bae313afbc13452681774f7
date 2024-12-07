package object_metadata

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
)

var thePool interfaces.Pool[Metadata, *Metadata]

func init() {
	thePool = pool.MakePool[Metadata, *Metadata](
		nil,
		Resetter.Reset,
	)
}

func GetPool() interfaces.Pool[Metadata, *Metadata] {
	return thePool
}
