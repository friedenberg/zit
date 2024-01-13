package metadatei

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/pool"
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
