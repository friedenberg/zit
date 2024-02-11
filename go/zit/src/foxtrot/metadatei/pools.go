package metadatei

import (
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit-go/src/bravo/pool"
)

var thePool schnittstellen.Pool[Metadatei, *Metadatei]

func init() {
	thePool = pool.MakePool[Metadatei, *Metadatei](
		nil,
		Resetter.Reset,
	)
}

func GetPool() schnittstellen.Pool[Metadatei, *Metadatei] {
	return thePool
}
