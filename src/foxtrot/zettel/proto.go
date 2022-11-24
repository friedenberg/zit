package zettel

import (
	"flag"

	"github.com/friedenberg/zit/src/charlie/bezeichnung"
	"github.com/friedenberg/zit/src/charlie/kennung"
	"github.com/friedenberg/zit/src/echo/typ"
)

type ProtoZettel struct {
	Typ         typ.Kennung
	Bezeichnung bezeichnung.Bezeichnung
	Etiketten   kennung.EtikettSet
}

func MakeProtoZettel() ProtoZettel {
	return ProtoZettel{
		//TODO-P2: use konfig to create correct default Typ
		Etiketten: kennung.MakeSet(),
	}
}

func (pz *ProtoZettel) AddToFlagSet(f *flag.FlagSet) {
	f.Var(&pz.Typ, "typ", "the Typ to use for created or updated Zettelen")
	f.Var(&pz.Bezeichnung, "bezeichnung", "the Bezeichnung to use for created or updated Zettelen")
	f.Var(&pz.Etiketten, "etiketten", "the Etiketten to use for created or updated Zttelen")
}

func (pz ProtoZettel) Equals(z Zettel) (ok bool) {
	var okTyp, okEt, okBez bool

	if !pz.Typ.IsEmpty() && pz.Typ.Equals(z.Typ) {
		okTyp = true
	}

	if pz.Etiketten.Len() > 0 && pz.Etiketten.Equals(z.Etiketten) {
		okEt = true
	}

	if !pz.Bezeichnung.WasSet() || pz.Bezeichnung.Equals(z.Bezeichnung) {
		okBez = true
	}

	ok = okTyp && okEt && okBez

	return
}

func (pz ProtoZettel) Make() (z *Zettel) {
	z = &Zettel{
		Etiketten: kennung.MakeSet(),
	}

	pz.Apply(z)

	return
}

func (pz ProtoZettel) Apply(z *Zettel) (ok bool) {
	if z.Typ.IsEmpty() && !pz.Typ.IsEmpty() && !z.Typ.Equals(pz.Typ) {
		ok = true
		z.Typ = pz.Typ
	}

	if pz.Bezeichnung.WasSet() && !z.Bezeichnung.Equals(pz.Bezeichnung) {
		ok = true
		z.Bezeichnung = pz.Bezeichnung
	}

	if pz.Etiketten.Len() > 0 {
		ok = true
	}

	mes := z.Etiketten.MutableCopy()
	pz.Etiketten.Each(mes.Add)
	z.Etiketten = mes.Copy()

	return
}
