package zettel

import (
	"github.com/friedenberg/zit/src/alfa/bezeichnung"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/typ"
)

type ProtoZettel struct {
	Typ         typ.Typ
	Bezeichnung *bezeichnung.Bezeichnung
	Etiketten   etikett.Set
}

func MakeProtoZettel() ProtoZettel {
	return ProtoZettel{
		Etiketten: etikett.MakeSet(),
	}
}

func (pz ProtoZettel) Equals(z Zettel) (ok bool) {
	var okTyp, okEt bool

	if !pz.Typ.IsEmpty() && pz.Typ.Equals(z.Typ) {
		okTyp = true
	}

	if pz.Etiketten.Len() > 0 && pz.Etiketten.Equals(z.Etiketten) {
		okEt = true
	}

	ok = okTyp && okEt

	return
}

func (pz ProtoZettel) Apply(z *Zettel) (ok bool) {
	if z.Typ.IsEmpty() && !pz.Typ.IsEmpty() && !z.Typ.Equals(pz.Typ) {
		ok = true
		z.Typ = pz.Typ
	}

	if pz.Bezeichnung != nil && !z.Bezeichnung.Equals(*pz.Bezeichnung) {
		ok = true
		z.Bezeichnung = *pz.Bezeichnung
	}

	if pz.Etiketten.Len() > 0 {
		ok = true
	}

	mes := z.Etiketten.MutableCopy()
	pz.Etiketten.Each(mes.Add)
	z.Etiketten = mes.Copy()

	return
}
