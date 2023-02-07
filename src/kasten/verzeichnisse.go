package kasten

import "github.com/friedenberg/zit/src/charlie/tridex"

type Verzeichnisse struct {
	wasPopulated bool
	Akten        *tridex.Tridex
	Objekten     *tridex.Tridex
}

func (z *Verzeichnisse) ResetWithObjekte(z1 Objekte) {
	z.wasPopulated = true
}

func (z *Verzeichnisse) Reset() {
	z.wasPopulated = false
	z.Akten = tridex.Make()
	z.Objekten = tridex.Make()
}

func (z *Verzeichnisse) ResetWith(z1 Verzeichnisse) {
	z.wasPopulated = true
	z.Akten = z1.Akten.Copy()
	z.Objekten = z1.Objekten.Copy()
}
