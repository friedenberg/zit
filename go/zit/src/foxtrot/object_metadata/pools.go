package object_metadata

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
)

var thePool interfaces.Pool[Metadatei, *Metadatei]

func init() {
	thePool = pool.MakePool[Metadatei, *Metadatei](
		nil,
		Resetter.Reset,
	)
}

func GetPool() interfaces.Pool[Metadatei, *Metadatei] {
	return thePool
}
