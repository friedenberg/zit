package kasten

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/tridex"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type Verzeichnisse struct {
	wasPopulated bool
	Akten        schnittstellen.MutableTridex
	Objekten     schnittstellen.MutableTridex
}

func (z *Verzeichnisse) ResetWithObjekteMetadateiGetter(
	z1 Objekte,
	_ metadatei.Getter,
) {
	z.wasPopulated = true
}

func (z *Verzeichnisse) Reset() {
	z.wasPopulated = false
	z.Akten = tridex.Make()
	z.Objekten = tridex.Make()
}

func (z *Verzeichnisse) ResetWith(z1 Verzeichnisse) {
	z.wasPopulated = true
	z.Akten = z1.Akten.MutableClone()
	z.Objekten = z1.Objekten.MutableClone()
}
