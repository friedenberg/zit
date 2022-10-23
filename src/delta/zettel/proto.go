package zettel

import (
	"github.com/friedenberg/zit/src/alfa/bezeichnung"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/typ"
)

type ProtoZettel struct {
	Typ         *typ.Typ
	Bezeichnung *bezeichnung.Bezeichnung
	Etiketten   etikett.Set
}

func (pz ProtoZettel) Apply(z *Zettel) {
	if pz.Typ != nil {
		z.Typ = *pz.Typ
	}

	if pz.Bezeichnung != nil {
		z.Bezeichnung = *pz.Bezeichnung
	}

	mes := z.Etiketten.MutableCopy()
	mes.Merge(pz.Etiketten)
	z.Etiketten = mes.Copy()
}
